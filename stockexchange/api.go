package stockexchange

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/svett/giraffe"
)

func Balance(w http.ResponseWriter, request *http.Request) {
	err := InSession(w, request, func(p *Portfolio) error {
		return nil
	})

	if err != nil {
		HTTPError(w, request, err.Error(), http.StatusInternalServerError)
	}
}

func Search(w http.ResponseWriter, request *http.Request) {
	query := request.FormValue("query")
	if query == "" {
		HTTPError(w, request, "The 'query' parameter is missing", http.StatusBadRequest)
		return
	}
	client := &Client{URL: "http://data.benzinga.com/rest"}
	result, err := client.Search(query)
	if err != nil {
		code := http.StatusInternalServerError
		if IsNotExistSybmol(err) {
			code = http.StatusNotFound
		}
		HTTPError(w, request, err.Error(), code)
	}

	encoder := giraffe.NewHTTPEncoder(w)
	encoder.EncodeJSON(result)
}

func Buy(w http.ResponseWriter, request *http.Request) {
	quantityParam := request.FormValue("quantity")
	if quantityParam == "" {
		HTTPError(w, request, "The 'quantity' parameter is missing", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityParam)
	if err != nil {
		HTTPError(w, request, err.Error(), http.StatusBadRequest)
		return
	}

	var stock Stock

	if err := json.NewDecoder(request.Body).Decode(&stock); err != nil {
		HTTPError(w, request, err.Error(), http.StatusBadRequest)
		return
	}

	err = InSession(w, request, func(p *Portfolio) error {
		return p.Buy(&stock, quantity)
	})

	if err != nil {
		HTTPError(w, request, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Sell(w http.ResponseWriter, request *http.Request) {
	var invoice Invoice

	if err := json.NewDecoder(request.Body).Decode(&invoice); err != nil {
		HTTPError(w, request, err.Error(), http.StatusBadRequest)
		return
	}

	err := InSession(w, request, func(p *Portfolio) error {
		return p.Sell(&invoice)
	})

	if err != nil {
		HTTPError(w, request, err.Error(), http.StatusInternalServerError)
	}
}

func HTTPError(w http.ResponseWriter, r *http.Request, err string, code int) {
	http.Error(w, err, code)
	log.Printf("Request: %s Method: %s Error: %s", r.URL.String(), r.Method, err)
}
