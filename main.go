package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/pawlobanano/csv-reader/customerimporter"
)

var log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	config, err := customerimporter.LoadConfig(log, ".env")
	if err != nil {
		log.Error("Loading config.", err)
		return
	}

	start := time.Now()

	err = customerimporter.Run(log, config)
	if err != nil {
		log.Error("Run customerimporter package.", err)
		return
	}

	log.Info("Program finished.", slog.String("time_taken", time.Since(start).String()))
}
