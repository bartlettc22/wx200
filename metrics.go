package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricSamplesCollected = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_samples_collected_total",
		Help: "Total number of data samples collected",
	}, []string{})
	metricChecksumFailures = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_checksum_failures_total",
		Help: "Total number of data checksum failures",
	}, []string{})
	metricHumidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_humidity",
		Help: "Current humidity measured as a percentage between 0-1",
	}, []string{"location"})
	metricGustSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_gust_speed",
		Help: "Current gust speeds measured in m/s",
	}, []string{})
	metricGustDir = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_gust_direction",
		Help: "Current gust direction measured in degrees",
	}, []string{})
	metricAvgWindSpeed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_avg_wind_speed",
		Help: "Average wind speed measured in m/s",
	}, []string{})
	metricAvgWindDir = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_avg_wind_direction",
		Help: "Average wind direction measured in degrees",
	}, []string{})
)
