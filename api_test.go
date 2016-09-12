package stockexchange_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/svett/stockexchange"
)

var _ = Describe("API", func() {
	var (
		handler     http.Handler
		cookie      *cookiejar.Jar
		client      *http.Client
		server      *httptest.Server
		stock       *stockexchange.Stock
		stockJSON   io.Reader
		invoice     *stockexchange.Invoice
		invoiceJSON io.Reader
	)

	BeforeEach(func() {
		cookie, _ = cookiejar.New(nil)
		client = &http.Client{
			Jar: cookie,
		}

		stock = &stockexchange.Stock{
			Symbol:   "B",
			Name:     "Bengaza",
			AskPrice: 5,
			BidPrice: 10,
		}

		invoice = &stockexchange.Invoice{
			Symbol:   "B",
			Price:    5,
			Quantity: 3,
		}

		stockData, err := json.Marshal(stock)
		Expect(err).NotTo(HaveOccurred())
		stockJSON = bytes.NewBuffer(stockData)

		invoiceData, err := json.Marshal(invoice)
		Expect(err).NotTo(HaveOccurred())
		invoiceJSON = bytes.NewBuffer(invoiceData)
	})

	JustBeforeEach(func() {
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Balance", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(stockexchange.Balance)
		})

		It("returns the portfolio", func() {
			resp, err := client.Get(server.URL)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result stockexchange.Portfolio
			Expect(json.NewDecoder(resp.Body).Decode(&result)).To(Succeed())
		})
	})

	Describe("Search", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(stockexchange.Search)
		})

		Context("when query parameter is missing", func() {
			It("returns an error", func() {
				resp, err := client.Get(server.URL)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				data, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("The 'query' parameter is missing\n"))
			})
		})

		It("returns a JSON output", func() {
			resp, err := client.Get(fmt.Sprintf("%s/search?query=A", server.URL))
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result []stockexchange.Stock
			Expect(json.NewDecoder(resp.Body).Decode(&result)).To(Succeed())
		})

		Context("when the symbol does not exists", func() {
			It("returns status code 404", func() {
				resp, err := client.Get(fmt.Sprintf("%s/search?query=J", server.URL))
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("Buy", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(stockexchange.Buy)
		})

		It("buys a stock", func() {
			resp, err := client.Post(fmt.Sprintf("%s/buy?quantity=2", server.URL), "application/json", stockJSON)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a JSON output", func() {
			resp, err := client.Post(fmt.Sprintf("%s/buy?quantity=2", server.URL), "application/json", stockJSON)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result stockexchange.Portfolio
			Expect(json.NewDecoder(resp.Body).Decode(&result)).To(Succeed())
		})

		Context("when stock is not provided", func() {
			It("returns an error", func() {
				resp, err := client.Post(server.URL, "application/json", nil)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the quantity is not provided", func() {
			It("returns an error", func() {
				resp, err := client.Post(server.URL, "application/json", stockJSON)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				data, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("The 'quantity' parameter is missing\n"))
			})
		})

		Context("when the quantity is not integer", func() {
			It("returns an error", func() {
				resp, err := client.Post(fmt.Sprintf("%s/buy?quantity=why", server.URL), "application/json", stockJSON)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the purchase operation fails", func() {
			It("returns an error", func() {
				resp, err := client.Post(fmt.Sprintf("%s/buy?quantity=-1", server.URL), "application/json", stockJSON)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Sell", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/buy") {
					stockexchange.Buy(w, r)
				} else if strings.HasPrefix(r.URL.Path, "/sell") {
					stockexchange.Sell(w, r)
				}
			})
		})

		JustBeforeEach(func() {
			resp, err := client.Post(fmt.Sprintf("%s/buy?quantity=2", server.URL), "application/json", stockJSON)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		// Due to limitations of client we cannot preserve the session
		XIt("sells shares", func() {
			resp, err := client.Post(fmt.Sprintf("%s/sell", server.URL), "application/json", invoiceJSON)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		// Due to limitations of client we cannot preserve the session
		XIt("returns a JSON output", func() {
			resp, err := client.Post(fmt.Sprintf("%s/sell", server.URL), "application/json", invoiceJSON)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result stockexchange.Portfolio
			Expect(json.NewDecoder(resp.Body).Decode(&result)).To(Succeed())
		})

		Context("when the sell operation fails", func() {
			It("returns the error", func() {
				resp, err := client.Post(fmt.Sprintf("%s/sell", server.URL), "application/json", nil)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
