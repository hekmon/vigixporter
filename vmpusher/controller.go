package vmpusher

// New returns an initialized and ready to use Controller
func New(levels, flows map[string]JSONLineMetric) *Controller {
	if levels == nil {
		levels = make(map[string]JSONLineMetric)
	}
	if flows == nil {
		flows = make(map[string]JSONLineMetric)
	}
	return &Controller{
		levels: levels,
		flows:  flows,
	}
}

// Controller handles the communication with the victoria metrics server
type Controller struct {
	levels map[string]JSONLineMetric
	flows  map[string]JSONLineMetric
}
