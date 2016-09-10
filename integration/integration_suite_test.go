package integration_test

import (
	"time"

	"github.com/svett/stockexchange/integration/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var runner *utils.StockExchangeRunner

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	runner = &utils.StockExchangeRunner{
		Stdout:            GinkgoWriter,
		Stderr:            GinkgoWriter,
		StartCheck:        "StackExchange started.",
		StartCheckTimeout: 3 * time.Second,
	}

	Expect(runner.Compile()).To(Succeed())
})

var _ = AfterSuite(func() {
	runner.Cleanup()
})
