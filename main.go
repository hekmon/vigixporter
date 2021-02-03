package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/watcher"
	systemd "github.com/iguanesolutions/go-systemd/v5"
	sysdnotify "github.com/iguanesolutions/go-systemd/v5/notify"

	"github.com/hekmon/hllogger"
)

var (
	logger        *hllogger.HlLogger
	core          *watcher.Controller
	mainLock      chan struct{}
	mainCtx       context.Context
	mainCtxCancel func()
)

func main() {
	listOfStations := []string{hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil}

	// Setup the logger
	_, systemdStarted := systemd.GetInvocationID()
	var logFlags int
	if !systemdStarted {
		logFlags = hllogger.LstdFlags
	}
	logger = hllogger.New(os.Stderr, &hllogger.Config{
		LogLevel:              hllogger.Debug,
		LoggerFlags:           logFlags,
		SystemdJournaldCompat: systemdStarted,
	})

	// Prepare main context for broadcasting the stop signal
	mainCtx, mainCtxCancel = context.WithCancel(context.Background())

	// Start core
	var err error
	if core, err = watcher.New(mainCtx, watcher.Config{
		Stations: listOfStations,
		Logger:   logger,
	}); err != nil {
		logger.Fatalf(1, "[Main] failed to instanciate the watcher: %s", err)
	}
	logger.Info("[Main] watcher started")

	// Everything is ready, listen to signals to know when to stop
	mainLock = make(chan struct{})
	go handleSignals()

	// Signal systemd we are ready if needed
	if err = sysdnotify.Ready(); err != nil {
		logger.Errorf("[Main] failed to notify systemd with ready signal: %s", err)
	}

	// Let's go to sleep while others do their work
	<-mainLock
}

func handleSignals() {
	var (
		sig os.Signal
		err error
	)
	// If we exit, allow main goroutine to do so
	defer close(mainLock)
	// Register signals
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)
	// Waiting for signals to catch
	for {
		sig = <-signalChannel
		switch sig {
		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGINT:
			logger.Infof("[Main] signal '%v' caught: cleaning up before exiting", sig)
			if err = sysdnotify.Stopping(); err != nil {
				logger.Errorf("[Main] can't send systemd stopping notification: %v", err)
			}
			// Cancel main ctx & wait for core to end
			mainCtxCancel()
			core.WaitStopped()
			logger.Debugf("[Main] signal '%v' caught: watcher stopped: unlocking main goroutine to exit", sig)
			return
		default:
			logger.Warningf("[Main] signal '%v' caught but no process set to handle it: skipping", sig)
		}
	}
}
