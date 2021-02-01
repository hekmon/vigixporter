package watcher

import (
	"context"
	"sync"
	"time"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/vmpusher"

	"github.com/hekmon/hllogger"
)

// Config allows to customize the instanciation of a watcher with New()
type Config struct {
	Stations []string

	Logger *hllogger.HlLogger
}

// New returns an initialized and ready to use Controller
func New(ctx context.Context, conf Config) (c *Controller, err error) {
	// Load state (TODO)
	var (
		levelsBuffer map[string]vmpusher.JSONLineMetric
		flowsBuffer  map[string]vmpusher.JSONLineMetric
		lastSeen     time.Time
	)
	// Init
	c = &Controller{
		stations: conf.Stations,
		lastSeen: lastSeen,
		logger:   conf.Logger,
		source:   hubeau.New(),
		target:   vmpusher.New(levelsBuffer, flowsBuffer),
		ctx:      ctx,
	}
	// Launch worker
	c.workers.Add(1)
	go func() {
		c.worker()
		c.workers.Done()
	}()
	// Create the auto-stopper (must be launch after the worker(s) in case ctx is cancelled while launching workers)
	go c.autostop()
	// Good to Go
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
	ctx     context.Context
	workers sync.WaitGroup
	stopped chan struct{}
}

func (c *Controller) autostop() {
	// Wait for signal
	<-c.ctx.Done()
	// Begin the stopping proceedure
	c.workers.Wait()
	// TODO: save some state ?
	// Close the stopped chan to indicate we are fully stopped
	close(c.stopped)
}

// WaitStopped will block until c is fully stopped.
// To be stopped, c needs to have its context cancelled.
// WaitStopped is safe to be called from multiples goroutines.
func (c *Controller) WaitStopped() {
	<-c.stopped
}
