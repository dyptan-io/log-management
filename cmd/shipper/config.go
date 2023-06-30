package main

import (
	"flag"
)

// Config is a struct that contains server configuration.
type Config struct {
	ReceiverAddr string
	WatchDir     string
}

func readConfig() Config {
	var config Config

	flag.StringVar(&config.ReceiverAddr, "receiver-addr", "http://localhost:80", "an address of receiver server")
	flag.StringVar(&config.WatchDir, "watch-dir", "./input", "a directory to watch for log files")
	flag.Parse()

	return config
}
