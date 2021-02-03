package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/watcher"

	"github.com/hekmon/hllogger"
	systemd "github.com/iguanesolutions/go-systemd/v5"
	sysdnotify "github.com/iguanesolutions/go-systemd/v5/notify"
)

const (
	confEnvarStations = "VIGIXPORTER_STATIONS"
	confEnvarLogLevel = "VIGIXPORTER_LOGLEVEL"
)

var (
	logger        *hllogger.HlLogger
	core          *watcher.Controller
	mainLock      chan struct{}
	mainCtx       context.Context
	mainCtxCancel func()
)

func main() {
	// Parse flags
	logLevelFlag := flag.String("loglevel", "info", "Set loglevel: debug, info, warning, error, fatal. Default info.")
	flag.Parse()

	// Init logger
	var logLevel hllogger.LogLevel
	switch strings.ToLower(*logLevelFlag) {
	case "debug":
		logLevel = hllogger.Debug
	case "info":
		logLevel = hllogger.Info
	case "warning":
		logLevel = hllogger.Warning
	case "error":
		logLevel = hllogger.Error
	case "fatal":
		logLevel = hllogger.Fatal
	default:
		logLevel = hllogger.Info
	}
	_, systemdStarted := systemd.GetInvocationID()
	var logFlags int
	if !systemdStarted {
		logFlags = hllogger.LstdFlags
	}
	logger = hllogger.New(os.Stderr, &hllogger.Config{
		LogLevel:              logLevel,
		LoggerFlags:           logFlags,
		SystemdJournaldCompat: systemdStarted,
	})

	// Get stations to follow from env
	var (
		stationsraw string
		stations    []string
	)
	if stationsraw = os.Getenv(confEnvarStations); stationsraw == "" {
		logger.Fatalf(1, "[Main] no stations submitted: use '%s' env var to set the stations to track. For example to follow Paris, Alfortville and Créteil: %s='%s,%s,%s'",
			confEnvarStations, confEnvarStations, hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil)
	}
	stations = strings.Split(stationsraw, ",")
	logger.Infof("[Main] conf: %d station(s) declared: %s", len(stations), strings.Join(stations, ", "))

	// Prepare main context for broadcasting the stop signal
	mainCtx, mainCtxCancel = context.WithCancel(context.Background())

	// Launch the watcher
	var err error
	if core, err = watcher.New(mainCtx, watcher.Config{
		Stations: stations,
		Logger:   logger,
	}); err != nil {
		logger.Fatalf(2, "[Main] failed to instanciate the watcher: %s", err)
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
