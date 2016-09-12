package stockexchange

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/svett/giraffe"
)

// store stores the portfolio as a cookie
var store = NewStore()

func init() {
	gob.Register(&Portfolio{})
}

// Store stores portfolios into a cookies
type Store struct {
	db *sessions.CookieStore
}

// NewStore creates a new store
func NewStore() *Store {
	return &Store{
		db: sessions.NewCookieStore([]byte("stock.exchange.store")),
	}
}

// Get a portfolio
func (store *Store) Get(r *http.Request) (*Portfolio, error) {
	session, err := store.session(r)
	if err != nil {
		return nil, err
	}

	value := session.Values["portfolio"]
	portfolio, ok := value.(*Portfolio)

	if !ok {
		portfolio = &Portfolio{
			Balance: 100000,
			Shares:  []*Share{},
		}
	}

	return portfolio, nil
}

// Save a portfolio
func (store *Store) Save(w http.ResponseWriter, r *http.Request, portfolio *Portfolio) error {
	session, err := store.session(r)
	if err != nil {
		return err
	}
	session.Values["portfolio"] = portfolio
	return session.Save(r, w)
}

func (store *Store) session(r *http.Request) (*sessions.Session, error) {
	return store.db.Get(r, "stock.exchange.session")
}

// InSession helper function that opens a session
func InSession(w http.ResponseWriter, request *http.Request, edit func(*Portfolio) error) error {
	p, err := store.Get(request)
	if err != nil {
		return err
	}

	if err := edit(p); err != nil {
		return err
	}

	if err := store.Save(w, request, p); err != nil {
		return err
	}

	if err := giraffe.NewHTTPEncoder(w).EncodeJSON(p); err != nil {
		return err
	}

	return nil
}
