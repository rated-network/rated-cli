package watcher

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
)

type Watcher struct {
	cfg *core.Config // Main configuration of rated CLI
}

func NewWatcher(cfg *core.Config) *Watcher {

	log.WithFields(log.Fields{
		"beacon-api-endpoint": cfg.BeaconEndpoint,
		"rated-api-endpoint": cfg.ApiEndpoint,
		"keys-to-watch":      len(cfg.WatcherValidationKeys),
		"refresh-rate":       cfg.WatcherRefreshRate,
	}).Info("created watcher")

	return &Watcher{
		cfg: cfg,
	}
}

func (w *Watcher) watchKeys() error {
	for _, key := range(w.cfg.WatcherValidationKeys) {
		log.WithFields(log.Fields{
			"validation-key": key,
		}).Info("fetching statistics about key")

		
	}

	return nil
}

func (w *Watcher) Watch() error {
	log.Info("starting to watch keys")

	for {
		startAt := time.Now()
		nextAt := startAt.Add(time.Duration(w.cfg.WatcherRefreshRate))
		log.WithFields(log.Fields{
			"start-at": startAt,
			"next-at": nextAt,
		}).Info("starting new iteration")

		err := w.watchKeys()
		if err != nil {
			log.WithError(err).Error("unable to watch keys")
			return err
		}

		sleepFor := time.Until(nextAt)
		log.WithFields(log.Fields{
			"sleep-for": sleepFor,
		}).Info("sleeping until next iteration")
		time.Sleep(sleepFor)
		log.Info("end of iteration")
	}

	return nil
}
