package main

import (
	"flag"
)

// Config contains server configuration
type Config struct {
	HTTPAddr string
}

// readConfig reads config values from command args
func readConfig() Config {
	config := Config{}

	flag.StringVar(&config.HTTPAddr, "addr", ":8080", "an address for HTTP server listener")
	flag.Parse()

	return config
}
