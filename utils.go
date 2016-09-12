package stockexchange

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// store is the backend session store
var store = sessions.NewCookieStore([]byte("this-should-be-an-env-variable"))

func init() {
	gob.Register(&Portfolio{})
}

func OpenSession(w http.ResponseWriter, r *http.Request, edit func(*Portfolio) error) error {
	session, err := store.Get(r, "stockbroker")
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

func HTTPError(w http.ResponseWriter, r *http.Request, err string, code int) {
	http.Error(w, err, code)
	log.Printf("Request: %s Method: %s Error: %s", r.URL.String(), r.Method, err)
}
