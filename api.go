package stockexchange

import (
	"net/http"

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
