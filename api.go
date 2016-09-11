package stockexchange

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/svett/giraffe"
)

// store is the backend session store
var store = sessions.NewFilesystemStore("", []byte("this-should-be-an-env-variable"))

func init() {
	gob.Register(&Portfolio{})
}

func Balance(w http.ResponseWriter, request *http.Request) {
	err := OpenSession(w, request, func(p *Portfolio) error {
		encoder := giraffe.NewHTTPEncoder(w)
		encoder.EncodeJSON(p)
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

	err = OpenSession(w, request, func(p *Portfolio) error {
		if err := p.Buy(&stock, quantity); err != nil {
			return err
		}
		encoder := giraffe.NewHTTPEncoder(w)
		encoder.EncodeJSON(p)
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Sell(w http.ResponseWriter, request *http.Request) {
	symbolParam := request.FormValue("symbol")
	if symbolParam == "" {
		http.Error(w, "The 'symbol' parameter is missing", http.StatusBadRequest)
		return
	}

	quantityParam := request.FormValue("quantity")
	if quantityParam == "" {
		http.Error(w, "The 'quantity' parameter is missing", http.StatusBadRequest)
		return
	}

	priceParam := request.FormValue("price")
	if priceParam == "" {
		http.Error(w, "The 'price' parameter is missing", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityParam)
	if err != nil {
		http.Error(w, "The provided quantity is not integer type", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceParam, 32)
	if err != nil {
		http.Error(w, "The provided price is not a valid numeric type", http.StatusBadRequest)
		return
	}

	err = OpenSession(w, request, func(p *Portfolio) error {
		if err := p.Sell(symbolParam, float32(price), quantity); err != nil {
			return err
		}

		encoder := giraffe.NewHTTPEncoder(w)
		encoder.EncodeJSON(p)
		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func OpenSession(w http.ResponseWriter, r *http.Request, edit func(*Portfolio) error) error {
	session, err := store.Get(r, "stockexchange")
	if err != nil {
		return err
	}

	value := session.Values["portfolio"]
	ok := true
	portfolio := &Portfolio{}

	if portfolio, ok = value.(*Portfolio); !ok {
		fmt.Println("NOT FOUND")
		portfolio = &Portfolio{
			Balance: 100000,
			Shares:  []*Share{},
		}
		session.Values["portfolio"] = portfolio
	}

	if err = session.Save(r, w); err != nil {
		return err
	}

	return edit(portfolio)
}
