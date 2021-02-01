package hubeau

import (
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

// New returns an initialized and ready to use Controller
func New() *Controller {
	return &Controller{
		http: cleanhttp.DefaultClient(),
	}
}

// Controller handles the communication with hubeau.
// Instanciate with New().
type Controller struct {
	http *http.Client
}
