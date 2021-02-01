package vmpusher

// New returns an initialized and ready to use Controller
func New() *Controller {
	return &Controller{
		levels: make(map[string]JSONLineMetric),
		flows:  make(map[string]JSONLineMetric),
	}
}

// Controller handles the communication with the victoria metrics server
type Controller struct {
	levels map[string]JSONLineMetric
	flows  map[string]JSONLineMetric
}
