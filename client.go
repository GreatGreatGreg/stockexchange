package stockexchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Stock represents a single stock option
type Stock struct {
	// Name of the company
	Name string `json:"name"`
	// Description of that company
	Description string `json:"description"`
	// Sector of this business
	Sector string `json:"sector"`
	// Industry of that sector
	Industry string `json:"industry"`
	// Symbol that identifies the stock
	Symbol string `json:"symbol"`
	// AskPrice is the price that you can buy that a single share
	AskPrice float32 `json:"askPrice"`
	// BidPrice is the price that you can sell a singel share
	BidPrice float32 `json:"bidPrice"`
}

// SearchError that can occur during search
type SearchError struct {
	// Code of the error
	Code int `json:"code"`
	// Message of the error
	Message string `json:"message"`
}

// Error returns the message
func (err SearchError) Error() string {
	return err.Message
}

// SearchResult is returned by the search
type SearchResult struct {
	Stock
	SearchError `json:"error"`
}

// Client that is used to retrieved data from Bengaza
type Client struct {
	// Bengaza API
	URL string
}

// Search looks up for particular symbol
func (client *Client) Search(symbol string) ([]Stock, error) {
	var result []Stock
	if client.URL == "" {
		return result, fmt.Errorf("The client 'URL' is not configured")
	}

	resp, err := http.Get(fmt.Sprintf("%s/richquoteDelayed?symbols=%s", client.URL, symbol))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}
		return result, fmt.Errorf(string(data))
	}

	data := make(map[string]SearchResult)
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return result, fmt.Errorf("The data cannot be decoded as JSON")
	}

	fmt.Println(data)

	for key, item := range data {
		if key == "null" {
			return []Stock{}, item.SearchError
		}
		result = append(result, item.Stock)
	}

	return result, nil
}
