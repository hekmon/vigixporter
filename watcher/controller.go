package watcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hekmon/vigixporter/hubeau/hydrometrie"
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
	// Load state
	previousState, err := loadState()
	if err != nil {
		err = fmt.Errorf("can not restore previous state: %w", err)
		return
	}
	// Init
	c = &Controller{
		stations:       conf.Stations,
		lastSeenLevels: previousState.LastSeenLevels,
		lastSeenFlows:  previousState.LastSeenFlows,
		logger:         conf.Logger,
		source:         hydrometrie.New(),
		target:         vmpusher.New(previousState.LevelsBuffer, previousState.FlowsBuffer),
		ctx:            ctx,
		stopped:        make(chan struct{}),
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
	lastSeenFlows  map[string]time.Time
	lastSeenLevels map[string]time.Time
	// subcontrollers
	logger *hllogger.HlLogger
	source *hydrometrie.Controller
	target *vmpusher.Controller
	// workers mgmt
	ctx     context.Context
	workers sync.WaitGroup
	stopped chan struct{}
}

func (c *Controller) autostop() {
	// Wait for signal
	<-c.ctx.Done()
	c.logger.Infof("[Watcher] stop signal received")
	// Begin the stopping proceedure
	c.workers.Wait()
	// Save state
	c.logger.Infof("[Watcher] worker stopped, dumping state to disk...")
	if err := saveState(state{
		LevelsBuffer:   c.target.GetLevelsBuffer(),
		FlowsBuffer:    c.target.GetFlowsBuffer(),
		LastSeenLevels: c.lastSeenLevels,
		LastSeenFlows:  c.lastSeenFlows,
	}); err != nil {
		c.logger.Errorf("[Watcher] failed to save state to disk, data will be lost: %s", err)
	} else {
		c.logger.Infof("[Watcher] state successfully written to disk")
	}
	// Close the stopped chan to indicate we are fully stopped
	close(c.stopped)
}

// WaitStopped will block until c is fully stopped.
// To be stopped, c needs to have its context cancelled.
// WaitStopped is safe to be called from multiples goroutines.
func (c *Controller) WaitStopped() {
	<-c.stopped
}
