package watcher

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
)

// Watcher watches a list of validation keys periodically.
type Watcher struct {
	cfg     *core.Config            // Main configuration of rated CLI
	keys    []EthereumValidationKey // List of keys we monitor
	reg     prometheus.Registerer   // Registerer of Prometheus metrics
	metrics WatcherMetrics          // Prometheus metrics for a validation key
}

// Representation of an Ethereum key and its associated statistics.
type EthereumValidationKey struct {
	publicKey string // validation key in the "0x..." format
	index     int    // index of the validation key on the blockchain
}

// CleanupValidationKey sanitizies the given validation key.
func cleanupValidationKey(key string) string {
	var prefix string

	if len(key) == 96 && !strings.HasPrefix(key, "0x") {
		log.WithFields(log.Fields{
			"key": key,
		}).Info("adding '0x' prefix to validation key")
		prefix = "0x"
	}

	return prefix + key
}

// NewWatcher creates a new watcher for validation keys.
func NewWatcher(cfg *core.Config, reg prometheus.Registerer) (*Watcher, error) {
	keys := []EthereumValidationKey{}

	log.WithFields(log.Fields{
		"beacon-api-endpoint": cfg.BeaconEndpoint,
		"rated-api-endpoint":  cfg.ApiEndpoint,
		"keys-to-watch":       len(cfg.WatcherValidationKeys),
		"refresh-rate":        cfg.WatcherRefreshRate,
	}).Info("created watcher")

	// Here we convert the validation keys into indexes, as this is what is
	// supported by the Rated Network API.
	for _, key := range cfg.WatcherValidationKeys {
		key = cleanupValidationKey(key)

		index, err := getValidationIndex(cfg, key)
		if err != nil {
			continue
		}

		log.WithFields(log.Fields{
			"validation-key":   key,
			"validation-index": index,
		}).Info("fetched validation key for the given index")

		keys = append(keys, EthereumValidationKey{
			publicKey: key,
			index:     index,
		})
	}

	return &Watcher{
		cfg:  cfg,
		keys: keys,
	}, nil
}

// Watch continuously fetches statistics about validation keys on rated.network.
func (w *Watcher) Watch() error {
	log.Info("starting to watch keys")

	for {
		startAt := time.Now()
		nextAt := startAt.Add(time.Duration(w.cfg.WatcherRefreshRate))
		log.WithFields(log.Fields{
			"start-at":        startAt,
			"next-at":         nextAt,
			"validation-keys": len(w.keys),
		}).Info("starting new iteration")

		for _, key := range w.keys {
			log.WithFields(log.Fields{
				"validation-key":       key.publicKey,
				"validation-key-index": key.index,
			}).Info("fetching statistics about key")

			stats, err := getValidationStatistics(w.cfg, &key)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"validation-key":       key.publicKey,
					"validation-key-index": key.index,
				}).Warn("unable to fetch statistics about key, skipped")
				continue
			}

			log.WithFields(log.Fields{
				"validation-key":          key.publicKey,
				"validation-key-index":    key.index,
				"uptime":                  stats.Uptime,
				"avg-correctness":         stats.AvgCorrectness,
				"attester-effectiveness":  stats.AttesterEffectiveness,
				"proposer-effectiveness":  stats.ProposerEffectiveness,
				"validator-effectiveness": stats.ValidatorEffectiveness,
			}).Info("fetched statistics about key from rated network")
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
