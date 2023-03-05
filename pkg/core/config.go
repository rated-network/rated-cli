package core

// Config represents the configuration of rated CLI.
type Config struct {
	ApiEndpoint           string   // endpoint to Rated network API
	ApiAccessToken        string   // access token for Rated network API
	Network               string   // network to query on Rated network API
	ListenOn              string   // interface to listen on to expose statistics
	Granularity           string   // Whether to fetch daily or hourly aggregates
	WatcherValidationKeys []string // validation keys to watch
}
