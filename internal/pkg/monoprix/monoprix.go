/*
Package monoprix is a scrapper for the monoprix website.

In order to bypass the security of the monoprix site, this package uses slimerjs (headless firefox) and TOR.
This scrapper is able to request a new circuit from TOR when necessary.

The scrapper itself is very easy to use, just call the Start function to start the scrapping. However, it must be ensured that TOR is launched, with the control port enabled. It is also necessary to ensure that slimerjs is accessible in the PATH. To simplify the setup, a docker image is available.

The scraper writes the results in JSON format to the requested file.
*/
package monoprix

import (
	"encoding/json"
	"fmt"
	"github.com/googollee/go-socket.io"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// MonoprixScraper is used to store the state of the scrapper. In concrete terms, this represents an instance of a scrapper.
type MonoprixScraper struct {
	config   Config
	waiting  []*work
	assigned []*work
	srv      *http.Server
	mutex    *sync.Mutex
	init     bool
}

// NewMonoprixScraper allows you to create a MonoprixScraper. This function takes as parameter the configuration of the new instance.
func NewMonoprixScraper(config Config) MonoprixScraper {
	return MonoprixScraper{
		config: config,
	}
}

// Start starts the scraping.
func (ms *MonoprixScraper) Start() error {
	log.SetLevel(log.DebugLevel)

	ms.mutex = &sync.Mutex{}
	ms.assigned = make([]*work, ms.config.Workers)

	f, err := os.Create(ms.config.OutputPath)
	if err != nil {
		log.WithField("error", err).Fatal("unable to open the output file:")
	}
	defer f.Close()

	server, err := socketio.NewServer(nil)
	if err != nil {
		return err
	}

	ms.mutex.Lock()
	ms.assigned[0] = &work{
		Link: "https://www.monoprix.fr/courses-en-ligne",
		Kind: "categories",
	}
	ms.mutex.Unlock()

	server.On("connection", func(so socketio.Socket) {
		log.Info("a worker is now connected")

		sendMessageToWorker(so, "init", initMessage{
			Workers: ms.config.Workers,
		})

		ms.init = true
		ms.next(so)

		so.On("categories", func(msg string) {
			var categoriesResultMessage categoriesResultMessage
			err := json.Unmarshal([]byte(msg), &categoriesResultMessage)
			if err != nil {
				log.WithField("error", err).Fatal("unable to unmarshal message:")
			}

			ms.mutex.Lock()
			for _, category := range categoriesResultMessage.Categories {
				ms.waiting = append(ms.waiting, &work{
					Link: category,
					Kind: "productsLinks",
				})
				//break // TODO: remove this, it's a test
			}
			ms.assigned[categoriesResultMessage.Worker] = nil
			ms.mutex.Unlock()

			ms.next(so)
		})

		so.On("productsLinks", func(msg string) {
			var productsLinksResultMessage productsLinksResultMessage
			err := json.Unmarshal([]byte(msg), &productsLinksResultMessage)
			if err != nil {
				log.WithField("error", err).Fatal("unable to unmarshal message:")
			}

			ms.mutex.Lock()
			if productsLinksResultMessage.NextPage != "" {
				ms.waiting = append(ms.waiting, &work{ // before to get new works faster
					Link: productsLinksResultMessage.NextPage,
					Kind: "productsLinks",
				})
			}

			for _, productLink := range productsLinksResultMessage.ProductsLinks {
				ms.waiting = append(ms.waiting, &work{
					Link: productLink,
					Kind: "product",
				})
			}

			ms.assigned[productsLinksResultMessage.Worker] = nil
			ms.mutex.Unlock()

			ms.next(so)
		})

		so.On("product", func(msg string) {
			var productResultMessage productResultMessage
			err := json.Unmarshal([]byte(msg), &productResultMessage)
			if err != nil {
				log.WithField("error", err).Fatal("unable to unmarshal message:")
			}

			ms.mutex.Lock()
			ms.assigned[productResultMessage.Worker] = nil
			ms.mutex.Unlock()

			b, err := json.Marshal(productResultMessage.Product)
			if err != nil {
				log.WithField("error", err).Fatal("unable to marshal the product:")
			}
			_, err = f.Write(b)
			if err != nil {
				log.WithField("error", err).Fatal("unable to write the product to the output file:")
			}

			ms.next(so)
		})

		so.On("errormsg", func(msg string) {
			var errorMessage errorMessage
			err := json.Unmarshal([]byte(msg), &errorMessage)
			if err != nil {
				log.WithField("error", err).Fatal("unable to unmarshal message:")
			}

			log.WithField("error", errorMessage).Warn("error received from a worker:")

			if errorMessage.Error == "Access Denied" {
				log.Debug("sending newnym to TOR...")

				err = newNym(ms.config)
				if err != nil {
					log.WithField("error", err).Fatal("unable to send newnym to TOR:")
				}

				err = spawnWorker(ms.config)
				if err != nil {
					log.WithField("error", err).Fatal("unable to spawn a new worker:")
				}
			}
		})

		so.On("disconnection", func() {
			log.Info("a worker is now disconnected")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	err = spawnWorker(ms.config)
	if err != nil {
		log.WithField("error", err).Fatal("unable to spawn a new worker:")
	}

	ms.srv = &http.Server{Addr: ":" + strconv.Itoa(ms.config.OrchestratorPort)}

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir(ms.config.AssetsPath)))
	log.WithField("configuration", ms.config).Info("listening:")
	err = ms.srv.ListenAndServe()
	if err != nil && err.Error() == "http: Server closed" {
		return nil
	}
	return err
}

// next is called after each event. It manages the state of the scrapper, and distribute the work.
func (ms *MonoprixScraper) next(so socketio.Socket) error {
	ms.mutex.Lock()
	if len(ms.waiting) == 0 { // TODO: is this ok?
		log.Debug("ms.waiting empty")
		assignedWorkEmpty := true
		for _, assignedWork := range ms.assigned {
			if assignedWork != nil {
				log.Debug("ms.assigned not empty")
				assignedWorkEmpty = false
				break
			}
		}

		if assignedWorkEmpty {
			log.Debug("closing http server...")
			ms.mutex.Unlock()
			return ms.end(so)
		}
	}

	if ms.init {
		for worker, assignedWork := range ms.assigned {
			if assignedWork != nil {
				err := sendMessageToWorker(so, "work", commandMessage{
					Command: "extract",
					Worker:  worker,
					Link:    assignedWork.Link,
					Kind:    assignedWork.Kind,
				})
				if err != nil {
					return err
				}
			}
		}

		ms.init = false
	}

	for worker, assignedWork := range ms.assigned {
		if len(ms.waiting) == 0 {
			break
		}
		if assignedWork == nil {
			ms.assigned[worker] = ms.waiting[0]
			ms.waiting = ms.waiting[1:]

			err := sendMessageToWorker(so, "work", commandMessage{
				Command: "extract",
				Worker:  worker,
				Link:    ms.assigned[worker].Link,
				Kind:    ms.assigned[worker].Kind,
			})
			if err != nil {
				return err
			}
		}
	}
	ms.mutex.Unlock()

	return nil
}

// end is called by next and indicates the end of the scrapping.
func (ms *MonoprixScraper) end(so socketio.Socket) error {
	err = ms.srv.Shutdown(nil)
	if err != nil {
		return err
	}

	return nil
}

// spawnBrowser launches a slimerjs process (a headless firefox browser)
func spawnWorker(config Config) error {
	cmd := exec.Command("slimerjs", config.AssetsPath+"/worker.js", fmt.Sprintf("--proxy=%s:%d", config.TORIP, config.TORPort), "--proxy-type=socks5", "--headless", "--load-images=no")
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

// sendMessageToBrowser allows you to send a message to the browser via a socket connection.
func sendMessageToWorker(so socketio.Socket, msgType string, msgData interface{}) error {
	msg, err := json.Marshal(msgData)
	if err != nil {
		return err
	}

	log.WithField("type", msgType).WithField("msg", string(msg)).Debug("sending this message to the worker:")

	return so.Emit(msgType, string(msg))
}
