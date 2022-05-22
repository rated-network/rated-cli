package core

import (
	"time"
)

// Config represents the configuration of rated CLI.
type Config struct {
	ApiEndpoint           string        // endpoint to Rated network API
	BeaconEndpoint        string        // endpoint to a beacon (lighthouse, teku, ...)
	WatcherRefreshRate    time.Duration // refresh rate of the watcher
	WatcherValidationKeys []string      // validation keys to watch
}
