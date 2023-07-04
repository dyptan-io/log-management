package main

import (
	"flag"
)

// Config contains server configuration for Receiver service.
type Config struct {
	HTTPAddr string
}

func readConfig() Config {
	config := Config{}

	flag.StringVar(&config.HTTPAddr, "addr", ":8080", "an address for HTTP server listener")
	flag.Parse()

	return config
}
