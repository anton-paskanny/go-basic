package api

import "demo/bin/config"

// Client defines interaction with JSON.BIN API
type Client interface {
	// TODO: define methods when API details are available
}

// client implements the Client interface
type client struct {
	config *config.Config
}

// New creates a new API client with configuration
func New(cfg *config.Config) Client {
	return &client{
		config: cfg,
	}
}
