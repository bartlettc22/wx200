package main

import (
	"fmt"
	"github.com/bartlettc22/wx200/pkg/wx200"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	// "sync"
	"time"
)

//
var wx *wx200.WX200

// MetricsHandler is our http request handler
type MetricsHandler struct{}

func main() {

	// var startupWG sync.WaitGroup
	// startupWG.Add(1)
	listenPort := 9041

	// Begin reading in WX200 data
	wx = wx200.New(&wx200.Config{
		ComPortName: "/dev/ttyUSB0",
		// StartupWG:   startupWG,
	})
	go wx.Go()

	// Wait until we've got data before serving
	// This avoids publishing zeros to Prometheus and potentially messing up data on restarts
	fmt.Println("Waiting for all data to come in before starting metrics server...")
	for {
		if !wx.Humidity.LastDataRecieved.IsZero() {
			break
		}
		time.Sleep(1 * time.Second)
	}

	var metricsHandler MetricsHandler
	http.Handle("/metrics", metricsHandler)
	fmt.Printf("Listening on port %d\n", listenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil)
}

// ServeHTTP gets the latest WX200 serial data and serves up the Prometheus metrics
func (m MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Set our metrics
	metricHumidity.With(prometheus.Labels{"location": "indoor"}).Set(float64(wx.Humidity.Indoor) / 100)

	// Let promhttp serve up the metrics page
	promHandler := promhttp.Handler()
	promHandler.ServeHTTP(w, r)
}
