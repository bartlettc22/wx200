package main

import (
	"fmt"
	"github.com/bartlettc22/wx200/pkg/wx200"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var version = "default"
var wx *wx200.WX200
var errorChan chan error
var windDataChan chan wx200.Wind
var windChillDataChan chan wx200.WindChill
var humidityDataChan chan wx200.Humidity
var rainDataChan chan wx200.Rain
var barometerDataChan chan wx200.Barometer
var dewPointDataChan chan wx200.DewPoint
var infoDataChan chan wx200.Info
var generalDataChan chan wx200.General
var temperatureDataChan chan wx200.Temperature
var listenPort string
var comPortName string
var v string

func main() {
	rootCmd.Execute()
}

func run() {
	// log.SetLevel(log.InfoLevel)
	log.Infof("Starting WX200 exporter v%s", version)

	// initialize our channels
	windDataChan = make(chan wx200.Wind, 1)
	windChillDataChan = make(chan wx200.WindChill, 1)
	humidityDataChan = make(chan wx200.Humidity, 1)
	barometerDataChan = make(chan wx200.Barometer, 1)
	dewPointDataChan = make(chan wx200.DewPoint, 1)
	rainDataChan = make(chan wx200.Rain, 1)
	infoDataChan = make(chan wx200.Info, 1)
	generalDataChan = make(chan wx200.General, 1)
	temperatureDataChan = make(chan wx200.Temperature, 1)
	errorChan = make(chan error)

	wx = wx200.New(&wx200.Config{
		ComPortName:         comPortName,
		WindDataChan:        windDataChan,
		WindChillDataChan:   windChillDataChan,
		HumidityDataChan:    humidityDataChan,
		RainDataChan:        rainDataChan,
		BarometerDataChan:   barometerDataChan,
		DewPointDataChan:    dewPointDataChan,
		ErrorChan:           errorChan,
		InfoDataChan:        infoDataChan,
		GeneralDataChan:     generalDataChan,
		TemperatureDataChan: temperatureDataChan,
	})

	// Start our collectors
	go watchErrors()
	go collectTemperatureMetrics()
	go collectWindMetrics()
	go collectWindChillMetrics()
	go collectHumidityMetrics()
	go collectBarometerMetrics()
	go collectDewPointMetrics()
	go collectRainMetrics()
	go collectInfoMetrics()
	go collectGeneralMetrics()

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

	log.Infof("Listening on port %s\n", listenPort)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
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
		metricWindChill.With(prometheus.Labels{}).Set(float64(w.Chill))
	}
}

func collectRainMetrics() {
	for r := range rainDataChan {
		log.WithFields(log.Fields{
			"rain": fmt.Sprintf("%+v", r),
		}).Debug("Rain data received")
		metricRain.With(prometheus.Labels{}).Set(float64(r.Total))
	}
}

func collectBarometerMetrics() {
	for m := range barometerDataChan {
		log.WithFields(log.Fields{
			"barometer": fmt.Sprintf("%+v", m),
		}).Debug("Barometer data received")
		metricBarometer.With(prometheus.Labels{"location": "local"}).Set(float64(m.Local))
	}
}

func collectDewPointMetrics() {
	for d := range dewPointDataChan {
		log.WithFields(log.Fields{
			"dewpoint": fmt.Sprintf("%+v", d),
		}).Debug("Dew Point data received")
		metricDewPoint.With(prometheus.Labels{"location": "indoor"}).Set(float64(d.Indoor))
		metricDewPoint.With(prometheus.Labels{"location": "outdoor"}).Set(float64(d.Outdoor))
	}
}

func collectInfoMetrics() {
	for i := range infoDataChan {
		metricSamplesCollected.With(prometheus.Labels{}).Set(float64(i.SamplesRecieved))
		metricChecksumFailures.With(prometheus.Labels{}).Set(float64(i.ChecksumFailures))
	}
}

func collectTemperatureMetrics() {
	for m := range temperatureDataChan {
		log.WithFields(log.Fields{
			"temperature": fmt.Sprintf("%+v", m),
		}).Debug("Temperature data received")
		metricTemperature.With(prometheus.Labels{"location": "indoor"}).Set(float64(m.Indoor))
		metricTemperature.With(prometheus.Labels{"location": "outdoor"}).Set(float64(m.Outdoor))
	}
}

func collectGeneralMetrics() {
	for g := range generalDataChan {

		if g.PowerSourceDC {
			metricGeneralPowerSource.With(prometheus.Labels{"source": "dc"}).Set(1)
			metricGeneralPowerSource.With(prometheus.Labels{"source": "ac"}).Set(0)
		} else {
			metricGeneralPowerSource.With(prometheus.Labels{"source": "dc"}).Set(0)
			metricGeneralPowerSource.With(prometheus.Labels{"source": "ac"}).Set(1)
		}

		if g.LowPowerIndicator {
			metricGeneralLowPower.With(prometheus.Labels{}).Set(1)
		} else {
			metricGeneralLowPower.With(prometheus.Labels{}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_CLOCK {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "clock"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "clock"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_TEMP {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "temperature"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "temperature"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_HUMIDITY {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "humidity"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "humidity"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_DEW {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "dewpoint"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "dewpoint"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_BARO {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "barometer"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "barometer"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_WIND {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "wind"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "wind"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_WINDCHILL {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "windchill"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "windchill"}).Set(0)
		}

		if g.DisplaySelected == wx200.DISPLAY_SELECTED_RAIN {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "rain"}).Set(1)
		} else {
			metricGeneralDisplaySelected.With(prometheus.Labels{"display": "rain"}).Set(0)
		}

		if g.DisplaySubscreen == wx200.DISPLAY_SUB_FIRST {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "first"}).Set(1)
		} else {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "first"}).Set(0)
		}

		if g.DisplaySubscreen == wx200.DISPLAY_SUB_SECOND {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "second"}).Set(1)
		} else {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "second"}).Set(0)
		}

		if g.DisplaySubscreen == wx200.DISPLAY_SUB_THIRD {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "third"}).Set(1)
		} else {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "third"}).Set(0)
		}

		if g.DisplaySubscreen == wx200.DISPLAY_SUB_FOURTH {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "fourth"}).Set(1)
		} else {
			metricGeneralDisplaySub.With(prometheus.Labels{"subscreen": "fourth"}).Set(0)
		}

		if g.DisplayType == wx200.DISPLAY_TYPE_MAIN {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "main"}).Set(1)
		} else {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "main"}).Set(0)
		}

		if g.DisplayType == wx200.DISPLAY_TYPE_MEM {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "memory"}).Set(1)
		} else {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "memory"}).Set(0)
		}

		if g.DisplayType == wx200.DISPLAY_TYPE_ALARM_IN {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "alarm_inside"}).Set(1)
		} else {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "alarm_inside"}).Set(0)
		}

		if g.DisplayType == wx200.DISPLAY_TYPE_ALARM_OUT {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "alarm_outside"}).Set(1)
		} else {
			metricGeneralDisplayType.With(prometheus.Labels{"type": "alarm_outside"}).Set(0)
		}

	}
}

func watchErrors() {
	for err := range errorChan {
		log.Warnf("%v", err)
	}
}
