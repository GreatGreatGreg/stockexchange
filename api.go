package stockexchange

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/svett/giraffe"
)

func Search(w http.ResponseWriter, request *http.Request) {
	query := request.FormValue("query")
	if query == "" {
		http.Error(w, "The 'query' parameter is missing", http.StatusBadRequest)
		return
	}
	client := &Client{URL: "http://data.benzinga.com/rest"}
	result, err := client.Search(query)
	if err != nil {
		code := http.StatusInternalServerError
		if IsNotExistSybmol(err) {
			code = http.StatusNotFound
		}
		http.Error(w, err.Error(), code)
	}

	encoder := giraffe.NewHTTPEncoder(w)
	encoder.EncodeJSON(result)
}

func Buy(w http.ResponseWriter, request *http.Request) {
	quantityParam := request.FormValue("quantity")
	if quantityParam == "" {
		http.Error(w, "The 'quantity' parameter is missing", http.StatusBadRequest)
		return
	}

	portfolio := &Portfolio{
		Balance: 10000,
		Shares:  []*Share{},
	}

	quantity, err := strconv.Atoi(quantityParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var stock Stock

	if err := json.NewDecoder(request.Body).Decode(&stock); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := portfolio.Buy(&stock, quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := giraffe.NewHTTPEncoder(w)
	encoder.EncodeJSON(portfolio)
}

func Sell(w http.ResponseWriter, request *http.Request) {
}
