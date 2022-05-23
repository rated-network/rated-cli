package core

import (
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServe(cfg *Config) {
	log.WithFields(log.Fields{
		"listen-on": cfg.ListenOn,
	}).Info("starting HTTP server")

	// Healthz handler used in Kubernetes set-ups to automatically restart
	// the container in case something goes off.
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Info("rater-cli is alive")
		w.WriteHeader(http.StatusOK)
	})

	// Prometheus handler to expose metrics to prometheus.
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(cfg.ListenOn , nil)
	if err != nil {
		log.WithError(err).Fatal("unable to watch")
	}

	log.Fatal("for some reason the HTTP server exited")
}
