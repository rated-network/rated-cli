package watcher

import (
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics exposed by the watcher.
type WatcherMetrics struct {
	ratedValidationUptime                 prometheus.Gauge
	ratedValidationAvgCorrectness         prometheus.Gauge
	ratedValidationAttesterEffectiveness  prometheus.Gauge
	ratedValidationProposerEffectiveness  prometheus.Gauge
	ratedValidationValidatorEffectiveness prometheus.Gauge
}

// NewWatcherMetrics creates prometheus metrics for the watcher.
func NewWatcherMetrics(reg prometheus.Registerer) *WatcherMetrics {
	log.Info("creating Prometheus metrics for watcher")

	uptime := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_uptime",
		Help:      "Uptime of a validation key.",
	})

	correctness := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_correctness",
		Help:      "Average correctness of a validation key.",
	})

	attester := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_attester_effectiveness",
		Help:      "Attester effectiveness of a validation key.",
	})

	proposer := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_proposer_effectiveness",
		Help:      "Proposer effectiveness of a validation key.",
	})

	effectiveness := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_effectiveness",
		Help:      "Effectiveness of a validation key.",
	})

	return &WatcherMetrics{
		ratedValidationUptime:                 uptime,
		ratedValidationAvgCorrectness:         correctness,
		ratedValidationAttesterEffectiveness:  attester,
		ratedValidationProposerEffectiveness:  proposer,
		ratedValidationValidatorEffectiveness: effectiveness,
	}
}
