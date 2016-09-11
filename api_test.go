package stockexchange_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/stockexchange"
)

var _ = Describe("API", func() {
	var (
		handler http.Handler
		server  *httptest.Server
	)

	JustBeforeEach(func() {
		server = httptest.NewServer(handler)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Search", func() {
		BeforeEach(func() {
			handler = http.HandlerFunc(stockexchange.Search)
		})

		Context("when query parameter is missing", func() {
			It("returns an error", func() {
				resp, err := http.Get(server.URL)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				data, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(data)).To(Equal("The 'query' parameter is missing\n"))
			})
		})

		It("returns a JSON output", func() {
			resp, err := http.Get(fmt.Sprintf("%s/search?query=A", server.URL))
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result []stockexchange.Stock
			Expect(json.NewDecoder(resp.Body).Decode(&result)).To(Succeed())
		})

		Context("when the symbol does not exists", func() {
			It("returns status code 404", func() {
				resp, err := http.Get(fmt.Sprintf("%s/search?query=J", server.URL))
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("Buy", func() {
		var (
			stockOne     *stockexchange.Stock
			stockOneJSON io.Reader
		)

		BeforeEach(func() {
			stockOne = &stockexchange.Stock{
				Symbol:   "B",
				Name:     "Bengaza",
				AskPrice: 5,
				BidPrice: 10,
			}

			data, err := json.Marshal(stockOne)
			Expect(err).NotTo(HaveOccurred())
			stockOneJSON = bytes.NewBuffer(data)

			handler = http.HandlerFunc(stockexchange.Buy)
		})

		It("buys a stock", func() {
			resp, err := http.Post(fmt.Sprintf("%s/buy?quantity=2", server.URL), "application/json", stockOneJSON)
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("returns a JSON output", func() {
			resp, err := http.Post(fmt.Sprintf("%s/buy?quantity=2", server.URL), "application/json", stockOneJSON)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var result stockexchange.Portfolio
			Expect(json.NewDecoder(resp.Body).Decode(&result)).To(Succeed())
		})

		Context("when stock is not provided", func() {
			It("returns an error", func() {
				resp, err := http.Post(server.URL, "application/json", nil)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the quantity is not provided", func() {
			It("returns an error", func() {
				resp, err := http.Post(server.URL, "application/json", stockOneJSON)
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
				resp, err := http.Post(fmt.Sprintf("%s/buy?quantity=why", server.URL), "application/json", stockOneJSON)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the purchase operation fails", func() {
			It("returns an error", func() {
				resp, err := http.Post(fmt.Sprintf("%s/buy?quantity=-1", server.URL), "application/json", stockOneJSON)
				Expect(err).NotTo(HaveOccurred())

				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
