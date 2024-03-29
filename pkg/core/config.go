package core

import "time"

// Config represents the configuration of rated CLI.
type Config struct {
	ApiEndpoint           string              // endpoint to Rated network API
	ApiAccessToken        string              // access token for Rated network API
	Network               string              // network to query on Rated network API
	ListenOn              string              // interface to listen on to expose statistics
	Granularity           string              // Whether to fetch daily or hourly aggregates
	WatcherValidationKeys map[string][]string // validation keys to watch
}

func (c *Config) SleepDuration() time.Duration {
	if c.Granularity == "hour" {
		return time.Hour
	} else {
		return 24 * time.Hour
	}
}

func (c *Config) Window() string {
	return c.Granularity
}
