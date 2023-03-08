package core

import "time"

// Config represents the configuration of rated CLI.
type Config struct {
	ApiEndpoint           string   // endpoint to Rated network API
	ApiAccessToken        string   // access token for Rated network API
	Network               string   // network to query on Rated network API
	ListenOn              string   // interface to listen on to expose statistics
	Granularity           string   // Whether to fetch daily or hourly aggregates
	WatcherValidationKeys []string // validation keys to watch
}

func (c *Config) SleepDuration() time.Duration {
	// Prater granularity defaults to 1 day
	if c.Network == "mainnet" && c.Granularity == "hour" {
		return time.Hour
	} else {
		return 24 * time.Hour
	}
}

func (c* Config) Window() string {
	// Prater granularity defaults to 1 day
	if c.Network == "mainnet" {
		return c.Granularity
	} else {
		return "day"
	}
}