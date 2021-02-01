package watcher

import (
	"context"
	"time"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/vmpusher"
)

// Config allows to customize the instanciation of a watcher with New()
type Config struct {
	Stations []string
}

// New returns an initialized and ready to use Controller
func New(ctx context.Context, conf Config) (c *Controller, err error) {
	c = &Controller{
		stations: conf.Stations,
		source:   hubeau.New(),
		target:   vmpusher.New(nil, nil),
		ctx:      ctx,
	}
	return
}

// Controller interfaces the watcher. Must be instanciated with New()
type Controller struct {
	stations []string
	source   *hubeau.Controller
	lastSeen time.Time
	target   *vmpusher.Controller
	ctx      context.Context
}
