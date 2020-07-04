package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// General
	metricGeneralPowerSource = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_power_source",
		Help: "Indicates the source of power for the WX200 device",
	}, []string{"source"})
	metricGeneralLowPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_power_low",
		Help: "Whether or not the WX200 device is running on low DC power (needs batteries)",
	}, []string{})
	metricGeneralDisplaySelected = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_display_selected",
		Help: "Indicates which display is selected on the WX200 device",
	}, []string{"display"})
	metricGeneralDisplaySub = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_display_subscreen",
		Help: "Indicates which display subscreen is selected on the WX200 device",
	}, []string{"subscreen"})
	metricGeneralDisplayType = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_display_type",
		Help: "Indicates which display type is selected on the WX200 device",
	}, []string{"type"})

	// Temperature
	metricTemperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_temperature",
		Help: "Current temperature by location (C)",
	}, []string{"location"})

	// Barometer
	metricBarometer = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_pressure",
		Help: "Current pressure by location (mb)",
	}, []string{"location"})

	// Dew Point
	metricDewPoint = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_dew_point",
		Help: "Current dew point by location (c)",
	}, []string{"location"})

	// Rain
	metricRain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_rain_total",
		Help: "Total cumulative rainfall",
	}, []string{})

	// Info
	metricSamplesCollected = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_samples_collected_total",
		Help: "Total number of data samples collected",
	}, []string{})
	metricChecksumFailures = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_checksum_failures_total",
		Help: "Total number of data checksum failures",
	}, []string{})

	// Humidity
	metricHumidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_humidity",
		Help: "Current humidity",
	}, []string{"location"})

	// Wind
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

	// Wind Chill
	metricWindChill = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wx200_wind_chill",
		Help: "Wind chill (C)",
	}, []string{})
)
