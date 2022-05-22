package watcher

import (
	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
)

type Watcher struct {
	cfg *core.Config // Main configuration of rated CLI
}

func NewWatcher(cfg *core.Config) *Watcher {

	log.WithFields(log.Fields{
		"rated-api-endpoint": cfg.ApiEndpoint,
		"keys-to-watch":      len(cfg.WatcherValidationKeys),
		"refresh-rate":       cfg.WatcherRefreshRate,
	}).Info("created watcher")

	return &Watcher{
		cfg: cfg,
	}
}

func (w *Watcher) Watch() error {
	return nil
}
