package watcher

import (
	log "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics exposed by the watcher.
type WatcherMetrics struct {
	ratedValidationUptime                 *prometheus.GaugeVec
	ratedValidationAvgCorrectness         *prometheus.GaugeVec
	ratedValidationAttesterEffectiveness  *prometheus.GaugeVec
	ratedValidationProposerEffectiveness  *prometheus.GaugeVec
	ratedValidationValidatorEffectiveness *prometheus.GaugeVec
}

// NewWatcherMetrics creates prometheus metrics for the watcher.
func NewWatcherMetrics(reg prometheus.Registerer) *WatcherMetrics {
	log.Info("creating Prometheus metrics for watcher")

	uptime := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_uptime",
		Help:      "Uptime of a validation key.",
	}, []string{"pubkey"})

	correctness := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_correctness",
		Help:      "Average correctness of a validation key.",
	}, []string{"pubkey"})

	attester := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_attester_effectiveness",
		Help:      "Attester effectiveness of a validation key.",
	}, []string{"pubkey"})

	proposer := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_proposer_effectiveness",
		Help:      "Proposer effectiveness of a validation key.",
	}, []string{"pubkey"})

	effectiveness := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "rated",
		Subsystem: "sentinel",
		Name:      "validation_key_effectiveness",
		Help:      "Effectiveness of a validation key.",
	}, []string{"pubkey"})

	return &WatcherMetrics{
		ratedValidationUptime:                 uptime,
		ratedValidationAvgCorrectness:         correctness,
		ratedValidationAttesterEffectiveness:  attester,
		ratedValidationProposerEffectiveness:  proposer,
		ratedValidationValidatorEffectiveness: effectiveness,
	}
}
