package vmpusher

func New() *Controller {
	return &Controller{
		levels: make(map[string]jsonlMetric),
		flows:  make(map[string]jsonlMetric),
	}
}

type Controller struct {
	levels map[string]jsonlMetric
	flows  map[string]jsonlMetric
}
