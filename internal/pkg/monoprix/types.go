package monoprix

// product represents a Monoprix product.
type product struct {
	Title       string `json:"title"`
	SubTitle    string `json:"subTitle"`
	Brand       string `json:"brand"`
	Price       string `json:"price"`
	WeightPrice string `json:"weightPrice"`
	URL         string `json:"url"`
}

// Config contains all the settings necessary for a scrapper instance.
type Config struct {
	AssetsPath         string
	OrchestratorPort   int
	TORIP              string
	TORPort            int
	TORControlPort     int
	TORControlPassword string
	OutputPath         string
	Workers            int
}

// categoriesResultMessage represents the result of a categories job.
type categoriesResultMessage struct {
	Worker     int      `json:"worker"`
	Categories []string `json:"categories"`
}

// productsLinksResultMessage represents the result of productsLinks job.
type productsLinksResultMessage struct {
	Worker        int      `json:"worker"`
	ProductsLinks []string `json:"productsLinks"`
	NextPage      string   `json:"nextPage"`
}

// productResultMessage represents the result of a product job.
type productResultMessage struct {
	Worker  int     `json:"worker"`
	Product product `json:"product"`
}

// commandMessage represents a command sent to the browser.
type commandMessage struct {
	Command string `json:"command"`
	Link    string `json:"link"`
	Kind    string `json:"kind"`
	Worker  int    `json:"worker"`
}

// errorMessage represents an error sent by the browser.
type errorMessage struct {
	Error string `json:"error"`
}

// work represents a job.
type work struct {
	Link string
	Kind string
}

// initMessage is sent to the browser and tells it how many workers (i.e. pages) to launch. This concretely indicates how many pages will be processed in parallel.
type initMessage struct {
	Workers int `json:"workers"`
}
