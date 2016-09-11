package stockexchange_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/stockexchange"
)

var _ = Describe("Client", func() {
	var client *stockexchange.Client

	BeforeEach(func() {
		client = &stockexchange.Client{}
	})

	Describe("Search", func() {
		var server *httptest.Server

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.URL.Path).To(Equal("/richquoteDelayed"))

				symbols := r.FormValue("symbols")

				if symbols == "A" {
					output, err := os.Open("fixtures/stock.txt")
					Expect(err).NotTo(HaveOccurred())
					defer output.Close()
					io.Copy(w, output)
				} else if symbols == "wrong_symbols" {
					fmt.Fprintln(w, "This Symbol is Wrong")
				} else if symbols == "errored_symbols" {
					http.Error(w, "Something wrong happened :(", http.StatusInternalServerError)
				} else if symbols == "not_found" {
					output, err := os.Open("fixtures/notfound.txt")
					Expect(err).NotTo(HaveOccurred())
					defer output.Close()
					io.Copy(w, output)
				}
			}))

			client.URL = server.URL
		})

		AfterEach(func() {
			server.Close()
		})

		It("returns a list of stock options", func() {
			stock, err := client.Search("A")
			Expect(stock).To(HaveLen(1))
			Expect(stock[0].Symbol).To(Equal("A"))
			Expect(stock[0].Name).To(Equal("Agilent Technologies"))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the URL is not provided", func() {
			BeforeEach(func() {
				client.URL = ""
			})

			It("return an error", func() {
				stock, err := client.Search("symbol")
				Expect(stock).To(HaveLen(0))
				Expect(err).To(MatchError("The client 'URL' is not configured"))
			})
		})

		Context("when the response is malformed JSON", func() {
			It("return an error", func() {
				stock, err := client.Search("wrong_symbols")
				Expect(stock).To(HaveLen(0))
				Expect(err).To(MatchError("The data cannot be decoded as JSON"))
			})
		})

		Context("when the server return an error", func() {
			It("return an error", func() {
				stock, err := client.Search("errored_symbols")
				Expect(stock).To(HaveLen(0))
				Expect(err).To(MatchError("Something wrong happened :(\n"))
			})
		})

		Context("when the symbol is not found", func() {
			It("return an error", func() {
				stock, err := client.Search("not_found")
				Expect(stock).To(HaveLen(0))
				Expect(err).To(MatchError("Unknown symbol."))
			})
		})
	})
})
