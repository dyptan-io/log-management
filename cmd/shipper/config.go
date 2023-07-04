package main

import (
	"flag"
	"strings"
)

// Config is a struct that contains Shipper service configuration.
type Config struct {
	ReceiverAddr string
	WatchDirs    []string
}

func readConfig() Config {
	var config Config

	var watchDirs string

	flag.StringVar(&watchDirs, "watch-dirs", "./testdata", "directories to watch for log files")
	flag.StringVar(&config.ReceiverAddr, "receiver-addr", "http://localhost:8080", "an address of the receiver server")
	flag.Parse()

	config.WatchDirs = strings.Split(watchDirs, ",")

	return config
}
