package main

import (
	"context"
	"os"
	"time"

	"github.com/hekmon/vigixporter/hubeau"
	"github.com/hekmon/vigixporter/watcher"

	"github.com/hekmon/hllogger"
)

var (
	logger *hllogger.HlLogger
)

func main() {
	listOfStations := []string{hubeau.StationParis, hubeau.StationAlfortville, hubeau.StationCreteil}

	logger = hllogger.New(os.Stderr, &hllogger.Config{
		LogLevel:    hllogger.Debug,
		LoggerFlags: hllogger.LstdFlags,
	})

	ctx, ctxCancel := context.WithCancel(context.Background())
	_, _ = watcher.New(ctx, watcher.Config{
		Stations: listOfStations,
		Logger:   logger,
	})
	time.Sleep(10*time.Minute + 30*time.Second)
	ctxCancel()
	time.Sleep(30 * time.Second)
}
