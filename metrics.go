package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricHumidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_humidity",
		Help: "Current humidity measured as a percentage between 0-1",
	}, []string{"location"})
)
