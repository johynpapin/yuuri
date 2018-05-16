package monoprix

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

// newNym requests a new circuit to TOR.
func newNym(config Config) error {
	subProcess := exec.Command("nc", "localhost", strconv.Itoa(config.TORControlPort))

	stdin, err := subProcess.StdinPipe()
	if err != nil {
		return err
	}

	if err = subProcess.Start(); err != nil {
		return err
	}

	io.WriteString(stdin, fmt.Sprintf("authenticate \"%s\"\n", config.TORControlPassword))
	io.WriteString(stdin, "signal newnym\n")
	io.WriteString(stdin, "quit\n")

	return subProcess.Wait()
}
