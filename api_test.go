package stockexchange_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/stockexchange"
)

var _ = Describe("API", func() {
	var server *httptest.Server

	BeforeEach(func() {
		server = httptest.NewServer(http.HandlerFunc(stockexchange.Search))
	})

	AfterEach(func() {
		server.Close()
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
