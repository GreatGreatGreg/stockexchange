package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type StockExchangeRunner struct {
	// path to the actual binary
	path string
	// env envorionemtn variables
	env []string

	// Stderr standard output stream
	Stdout io.Writer
	// Stderr standard error stream
	Stderr io.Writer
	// StartCheckTimeout that timeout when the app is not started
	StartCheckTimeout time.Duration
	// StartCheck text to match to indicate sucessful start.
	StartCheck string
}

func (ex *StockExchangeRunner) Compile() error {
	bin, err := gexec.Build("github.com/svett/stockExchange/cmd/stockexchange")
	if err != nil {
		return err
	}

	ex.env = []string{}
	ex.path = bin
	return nil
}

func (ex *StockExchangeRunner) Setenv(name, value string) {
	ex.env = append(ex.env, fmt.Sprintf("%s=%s", name, value))
}

func (ex *StockExchangeRunner) Clearenv() {
	ex.env = []string{}
}

func (ex *StockExchangeRunner) Start(args ...string) (*gexec.Session, error) {
	buffer := gbytes.NewBuffer()
	stdout := io.MultiWriter(buffer, ex.Stdout)
	stderr := io.MultiWriter(buffer, ex.Stderr)

	cmd := exec.Command(ex.path, args...)
	cmd.Env = append(os.Environ(), ex.env...)
	session, err := gexec.Start(cmd, stdout, stderr)

	timeout := time.After(ex.StartCheckTimeout)
	detector := buffer.Detect(ex.StartCheck)

	for {
		select {
		case <-detector:
			buffer.CancelDetects()
			return session, err
		case <-timeout:
			session.Kill().Wait()
			return nil, fmt.Errorf("did not see %s in command's output within %s: %s", ex.StartCheck, ex.StartCheckTimeout, string(buffer.Contents()))
		}
	}
}

func (ex *StockExchangeRunner) Cleanup() {
	gexec.CleanupBuildArtifacts()
	ex.path = ""
}
