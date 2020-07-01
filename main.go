package main

import (
	"fmt"
	"github.com/bartlettc22/wx200/pkg/wx200"
	// "github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// Application version - passed in via build
var version = "default"

var wx *wx200.WX200
var errorChan chan error
var windDataChan chan wx200.Wind
var windChillDataChan chan wx200.WindChill
var humidityDataChan chan wx200.Humidity
var rainDataChan chan wx200.Rain
var infoDataChan chan wx200.Info

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {

	log.SetLevel(log.DebugLevel)
	log.Infof("Starting WX200 exporter v%s", version)

	// CMD VARS
	listenPort := 9041
	comPortName := "/dev/ttyUSB0"

	// initialize our channels
	windDataChan = make(chan wx200.Wind, 1)
	windChillDataChan = make(chan wx200.WindChill, 1)
	humidityDataChan = make(chan wx200.Humidity, 1)
	infoDataChan = make(chan wx200.Info, 1)
	errorChan = make(chan error)

	wx = wx200.New(&wx200.Config{
		ComPortName:       comPortName,
		WindDataChan:      windDataChan,
		WindChillDataChan: windChillDataChan,
		HumidityDataChan:  humidityDataChan,
		RainDataChan:      rainDataChan,
		ErrorChan:         errorChan,
		InfoDataChan:      infoDataChan,
	})

	// Start our collectors
	go watchErrors()
	go collectWindMetrics()
	go collectWindChillMetrics()
	go collectHumidityMetrics()
	go collectInfoMetrics()

	// Start serial communication
	go wx.Go()

	// Wait until we've got data before serving the metrics endpoint
	// This avoids publishing zeros to Prometheus and potentially messing up data on restarts
	log.Info("Waiting for all data to come in before starting metrics server...")
	for {
		if wx.Ready() {
			break
		}
		time.Sleep(1 * time.Second)
	}

	log.Infof("Listening on port %d\n", listenPort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

func collectHumidityMetrics() {
	for h := range humidityDataChan {
		log.WithFields(log.Fields{
			"humidity": fmt.Sprintf("%+v", h),
		}).Debug("Humidity received")
		metricHumidity.With(prometheus.Labels{"location": "indoor"}).Set(float64(h.Indoor) / 100)
		metricHumidity.With(prometheus.Labels{"location": "outdoor"}).Set(float64(h.Outdoor) / 100)
	}
}

func collectWindMetrics() {
	for w := range windDataChan {
		log.WithFields(log.Fields{
			"wind": fmt.Sprintf("%+v", w),
		}).Debug("Wind data received")
		metricGustSpeed.With(prometheus.Labels{}).Set(w.GustSpeed)
		metricGustDir.With(prometheus.Labels{}).Set(float64(w.GustDirection))
		metricAvgWindSpeed.With(prometheus.Labels{}).Set(w.AvgSpeed)
		metricAvgWindDir.With(prometheus.Labels{}).Set(float64(w.AvgDirection))
	}
}

func collectWindChillMetrics() {
	for w := range windChillDataChan {
		log.WithFields(log.Fields{
			"windchill": fmt.Sprintf("%+v", w),
		}).Debug("Wind chill data received")
		// metricGustSpeed.With(prometheus.Labels{}).Set(w.GustSpeed)
		// metricGustDir.With(prometheus.Labels{}).Set(float64(w.GustDirection))
		// metricAvgWindSpeed.With(prometheus.Labels{}).Set(w.AvgSpeed)
		// metricAvgWindDir.With(prometheus.Labels{}).Set(float64(w.AvgDirection))
	}
}

func collectRainMetrics() {
	for r := range rainDataChan {
		log.WithFields(log.Fields{
			"rain": fmt.Sprintf("%+v", r),
		}).Debug("Rain data received")
		// spew.Dump(r)
		// metricGustSpeed.With(prometheus.Labels{}).Set(w.GustSpeed)
		// metricGustDir.With(prometheus.Labels{}).Set(float64(w.GustDirection))
		// metricAvgWindSpeed.With(prometheus.Labels{}).Set(w.AvgSpeed)
		// metricAvgWindDir.With(prometheus.Labels{}).Set(float64(w.AvgDirection))
	}
}

func collectInfoMetrics() {
	for i := range infoDataChan {
		metricSamplesCollected.With(prometheus.Labels{}).Set(float64(i.SamplesRecieved))
		metricChecksumFailures.With(prometheus.Labels{}).Set(float64(i.ChecksumFailures))
	}
}

func watchErrors() {
	for err := range errorChan {
		log.Warnf("%v", err)
	}
}
