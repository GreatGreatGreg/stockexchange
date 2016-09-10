package stockexchange

import (
	"net/http"

	"github.com/svett/giraffe"
)

// HTMLPage
type HTMLPage struct {
	// Title
	Title string
}

func Index(w http.ResponseWriter, req *http.Request) {
	renderer := giraffe.NewHTMLTemplateRenderer(w)
	renderer.Render("index", HTMLPage{Title: "StackExchange"})
}
