package watcher

import (
	"context"
	"time"

	"github.com/hekmon/hllogger"
	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/vmpusher"
)

// Config allows to customize the instanciation of a watcher with New()
type Config struct {
	Stations     []string
	LevelsBuffer map[string]vmpusher.JSONLineMetric
	FlowsBuffer  map[string]vmpusher.JSONLineMetric
	Logger       *hllogger.HlLogger
}

// New returns an initialized and ready to use Controller
func New(ctx context.Context, conf Config) (c *Controller, err error) {
	c = &Controller{
		stations: conf.Stations,
		logger:   conf.Logger,
		source:   hubeau.New(),
		target:   vmpusher.New(conf.LevelsBuffer, conf.FlowsBuffer),
		ctx:      ctx,
	}
	return
}

// Controller interfaces the watcher. Must be instanciated with New()
type Controller struct {
	// config
	stations []string
	// state
	lastSeen time.Time
	// subcontrollers
	logger *hllogger.HlLogger
	source *hubeau.Controller
	target *vmpusher.Controller
	// workers mgmt
	ctx context.Context
}
