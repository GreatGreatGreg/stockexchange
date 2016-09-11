package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/onsi/gomega/gexec"
	"github.com/svett/stockexchange"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func() {
	var (
		args       []string
		session    *gexec.Session
		sessionErr error
	)

	JustBeforeEach(func() {
		session, sessionErr = runner.Start(args...)
	})

	AfterEach(func() {
		if session != nil {
			session.Kill()
		}
	})

	It("is listenting on HTTP port 9292", func() {
		Expect(sessionErr).NotTo(HaveOccurred())
		_, err := http.Get("http://127.0.0.1:9292")
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when command line argument is provided", func() {
		BeforeEach(func() {
			args = []string{"--addr=127.0.0.1:8080"}
		})

		It("does not overried the PORT", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			_, err := http.Get("http://127.0.0.1:8080")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("API", func() {
		var stockData io.Reader

		BeforeEach(func() {
			args = []string{"--addr=127.0.0.1:8080"}

			stock := &stockexchange.Stock{
				Symbol:   "B",
				Name:     "Bengaza",
				AskPrice: 5,
				BidPrice: 10,
			}

			data, err := json.Marshal(stock)
			Expect(err).NotTo(HaveOccurred())
			stockData = bytes.NewBuffer(data)
		})

		It("handles search requests", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			resp, err := http.Get("http://127.0.0.1:8080/api/v1/search?query=A")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("handles buy requests", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			resp, err := http.Post("http://127.0.0.1:8080/api/v1/buy?quantity=2", "application/json", stockData)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Describe("Environment variable PORT", func() {
		BeforeEach(func() {
			args = []string{}
			runner.Setenv("PORT", "8899")
		})

		AfterEach(func() {
			runner.Clearenv()
		})

		It("is listenting on that HTTP port", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			resp, err := http.Get("http://127.0.0.1:8899/api/v1/portfolio")
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("when the port is not integer", func() {
			BeforeEach(func() {
				os.Setenv("PORT", "wrong_port")
			})

			It("is listenting on the default HTTP port", func() {
				Expect(sessionErr).NotTo(HaveOccurred())
				_, err := http.Get("http://127.0.0.1:9292")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("when the addr is provided", func() {
		BeforeEach(func() {
			args = []string{"--addr=127.0.0.1:8080"}
		})

		It("is listenting on that HTTP port", func() {
			Expect(sessionErr).NotTo(HaveOccurred())
			_, err := http.Get("http://127.0.0.1:8080")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when the addr is wrong", func() {
		BeforeEach(func() {
			args = []string{"--addr=wrong_host_and_port"}
		})

		It("returns an error", func() {
			Expect(session).To(BeNil())
			Expect(sessionErr.Error()).To(ContainSubstring("The provided wrong_host_and_port addr is not correct"))
		})
	})
})
