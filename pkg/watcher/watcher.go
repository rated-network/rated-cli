package watcher

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/rated-cli/pkg/core"
)

// Watcher watches a list of validation keys periodically.
type Watcher struct {
	cfg     *core.Config          // Main configuration of rated CLI
	keys    []string              // List of keys we monitor
	reg     prometheus.Registerer // Registerer of Prometheus metrics
	metrics *WatcherMetrics       // Prometheus metrics for a validation key
}

// NewWatcher creates a new watcher for validation keys.
func NewWatcher(cfg *core.Config, reg prometheus.Registerer) (*Watcher, error) {
	log.WithFields(log.Fields{
		"rated-api-endpoint": cfg.ApiEndpoint,
		"keys-to-watch":      len(cfg.WatcherValidationKeys),
		"granularity":        cfg.Granularity,
	}).Info("created watcher")

	metrics := NewWatcherMetrics(reg)

	return &Watcher{
		cfg:     cfg,
		keys:    cfg.WatcherValidationKeys,
		reg:     reg,
		metrics: metrics,
	}, nil
}

// Watch continuously fetches statistics about validation keys on rated.network.
func (w *Watcher) Watch() error {
	log.Info("starting to watch keys")

	for {
		startAt := time.Now()
		granularity := w.cfg.Granularity
		sleepDuration := time.Duration(86400)
		if granularity == "hour" {
			sleepDuration = time.Duration(3600)
		}
		nextAt := startAt.Add(sleepDuration)
		log.WithFields(log.Fields{
			"start-at":        startAt,
			"next-at":         nextAt,
			"validation-keys": len(w.keys),
		}).Info("starting new iteration")

		w.metrics.ratedMonitoredKeys.Set(float64(len(w.keys)))

		for _, key := range w.keys {
			log.WithFields(log.Fields{
				"validation-key": key,
			}).Info("fetching statistics about key")

			stats, err := getValidationStatistics(w.cfg, key, granularity)
			if err != nil {
				log.WithError(err).WithFields(log.Fields{
					"validation-key": key,
				}).Warn("unable to fetch statistics about key, skipped")
				continue
			}

			w.metrics.ratedValidationUptime.WithLabelValues(key).Set(stats.Uptime)
			w.metrics.ratedValidationAvgCorrectness.WithLabelValues(key).Set(stats.AvgCorrectness)
			w.metrics.ratedValidationAttesterEffectiveness.WithLabelValues(key).Set(stats.AttesterEffectiveness)
			w.metrics.ratedValidationProposerEffectiveness.WithLabelValues(key).Set(stats.ProposerEffectiveness)
			w.metrics.ratedValidationValidatorEffectiveness.WithLabelValues(key).Set(stats.ValidatorEffectiveness)
			w.metrics.ratedValidationRewards.WithLabelValues(key).Set(stats.Rewards)
			w.metrics.ratedValidationInclusionDelay.WithLabelValues(key).Set(stats.InclusionDelay)

			log.WithFields(log.Fields{
				"validation-key":          key,
				"uptime":                  stats.Uptime,
				"avg-correctness":         stats.AvgCorrectness,
				"attester-effectiveness":  stats.AttesterEffectiveness,
				"proposer-effectiveness":  stats.ProposerEffectiveness,
				"validator-effectiveness": stats.ValidatorEffectiveness,
				"rewards":                 stats.Rewards,
				"inclusion-delay":         stats.InclusionDelay,
			}).Info("fetched statistics about key from rated network")
		}

		sleepFor := time.Until(nextAt)
		log.WithFields(log.Fields{
			"sleep-for": sleepFor,
		}).Info("sleeping until next iteration")
		time.Sleep(sleepFor)
		log.Info("end of iteration")
	}
}
