package rated

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Watcher struct {
	ratedApiEndpoint string        // endpoint to the rated API
	refreshRate      time.Duration // how often we should fetch statistics from keys
	validationKeys   []string      // keys to watch
}

func NewWatcher(ratedApiEndpoint string, validationKeys []string, refreshRateSec int64) *Watcher {
	refreshRate := time.Second * time.Duration(refreshRateSec)

	log.WithFields(log.Fields{
		"keys-to-watch": len(validationKeys),
		"rated-api-endpoint": ratedApiEndpoint,
		"refresh-rate": refreshRate,
	}).Info("created watcher")

	return &Watcher{
		ratedApiEndpoint: ratedApiEndpoint,
		refreshRate:      refreshRate,
		validationKeys:   validationKeys,
	}
}

func (w *Watcher) Watch() error {
	return nil
}
