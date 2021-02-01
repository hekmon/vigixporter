package vmpusher

func New() *Controller {
	return &Controller{
		levels: make(map[string]JSONLineMetric),
		flows:  make(map[string]JSONLineMetric),
	}
}

type Controller struct {
	levels map[string]JSONLineMetric
	flows  map[string]JSONLineMetric
}
