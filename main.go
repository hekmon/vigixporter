package main

import (
	"context"
	"os"
	"time"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/watcher"
	systemd "github.com/iguanesolutions/go-systemd/v5"

	"github.com/hekmon/hllogger"
)

var (
	logger *hllogger.HlLogger
)

func main() {
	listOfStations := []string{hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil}

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

	ctx, ctxCancel := context.WithCancel(context.Background())

	_, err := watcher.New(ctx, watcher.Config{
		Stations: listOfStations,
		Logger:   logger,
	})

	if err != nil {
		logger.Fatalf(1, "Failed to instanciate the watcher: %s", err)
	}

	time.Sleep(10*time.Minute + 30*time.Second)
	ctxCancel()
	time.Sleep(30 * time.Second)
}
