package vmpusher

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
)

// Config allows to configure a vmpusher controller with New()
type Config struct {
	VMURL   string
	VMUser  string
	VMPass  string
	Levels  map[string]JSONLineMetric
	Flows   map[string]JSONLineMetric
	Context context.Context
}

// New returns an initialized and ready to use Controller
func New(conf Config) *Controller {
	if conf.Levels == nil {
		conf.Levels = make(map[string]JSONLineMetric)
	}
	if conf.Flows == nil {
		conf.Flows = make(map[string]JSONLineMetric)
	}
	return &Controller{
		vmURL:  conf.VMURL,
		vmUser: conf.VMUser,
		vmPass: conf.VMPass,
		levels: conf.Levels,
		flows:  conf.Flows,
		http:   cleanhttp.DefaultClient(),
		ctx:    conf.Context,
	}
}

// Controller handles the communication with the victoria metrics server
type Controller struct {
	// conf
	vmURL  string
	vmUser string
	vmPass string
	// buffers
	levels map[string]JSONLineMetric
	flows  map[string]JSONLineMetric
	// sub controllers
	http *http.Client
	ctx  context.Context
}
