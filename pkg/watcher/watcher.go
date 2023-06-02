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
	keys    map[string][]string   // List of keys we monitor
	reg     prometheus.Registerer // Registerer of Prometheus metrics
	metrics *WatcherMetrics       // Prometheus metrics for a validation key
}

// NewWatcher creates a new watcher for validation keys.
func NewWatcher(cfg *core.Config, reg prometheus.Registerer) (*Watcher, error) {
	log.WithFields(log.Fields{
		"rated-api-endpoint": cfg.ApiEndpoint,
		"keys-to-watch":      countValidationKeys(cfg.WatcherValidationKeys),
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
		granularity := w.cfg.Window()
		nextAt := startAt.Add(w.cfg.SleepDuration())
		log.WithFields(log.Fields{
			"start-at":        startAt,
			"next-at":         nextAt,
			"validation-keys": countValidationKeys(w.cfg.WatcherValidationKeys),
		}).Info("starting new iteration")

		w.metrics.ratedMonitoredKeys.Set(float64(countValidationKeys(w.cfg.WatcherValidationKeys)))

		for label, keys := range w.cfg.WatcherValidationKeys {
			// Aggregate the number of validators per label
			validatorCountByLabel := AggregateValidatorsByLabel(map[string][]string{label: keys})

			// Update the monitored keys count metric per label
			for label, count := range validatorCountByLabel {
				w.metrics.ratedMonitoredByLabel.WithLabelValues(label).Set(float64(count))
			}
			for _, key := range keys {
				log.WithFields(log.Fields{
					"label":          label,
					"validation-key": key,
				}).Info("fetching statistics about key")

				stats, err := getValidationStatistics(w.cfg, key, granularity)
				if err != nil {
					log.WithError(err).WithFields(log.Fields{
						"label":          label,
						"validation-key": key,
					}).Warn("unable to fetch statistics about key, skipped")
					continue
				}

				metricLabels := prometheus.Labels{"label": label, "pubkey": key}
				w.metrics.ratedValidationUptime.With(metricLabels).Set(stats.Uptime)
				w.metrics.ratedValidationAvgCorrectness.With(metricLabels).Set(stats.AvgCorrectness)
				w.metrics.ratedValidationAttesterEffectiveness.With(metricLabels).Set(stats.AttesterEffectiveness)
				w.metrics.ratedValidationProposerEffectiveness.With(metricLabels).Set(stats.ProposerEffectiveness)
				w.metrics.ratedValidationValidatorEffectiveness.With(metricLabels).Set(stats.ValidatorEffectiveness)
				w.metrics.ratedValidationRewards.With(metricLabels).Set(stats.Rewards)
				w.metrics.ratedValidationInclusionDelay.With(metricLabels).Set(stats.InclusionDelay)

				log.WithFields(log.Fields{
					"label":                   label,
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
		}

		sleepFor := time.Until(nextAt)
		log.WithFields(log.Fields{
			"sleep-for": sleepFor,
		}).Info("sleeping until next iteration")
		time.Sleep(sleepFor)
		log.Info("end of iteration")
	}
}

// countValidationKeys counts the total number of validation keys in the configuration.
func countValidationKeys(validationKeys map[string][]string) int {
	count := 0
	for _, keys := range validationKeys {
		count += len(keys)
	}
	return count
}

// AggregateValidatorsByLabel aggregates the number of validators per label and returns a map with label and count.
func AggregateValidatorsByLabel(validationKeys map[string][]string) map[string]int {
	validatorCountByLabel := make(map[string]int)

	for label, keys := range validationKeys {
		validatorCountByLabel[label] = len(keys)
	}

	return validatorCountByLabel
}
