# Radio Shack WX-200 Electronic Weather Station
[![Build][Build-Status-Image]][Build-Status-Url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

This project contains the following:
* Prometheus exporter for WX200 readings (amd64/arm/Docker)
* Golang library for interfacing with the WX200 via serial communication


## Credits
This project could not have been completed without the work of Mike Wingstrom, Glynne Tolar and Tim Witham at [wx200.planetfall.com](http://wx200.planetfall.com/).  Their [serial protocol mapping](http://wx200.planetfall.com/wx200.txt) (also located [here](docs/wx200_serial_protocol.txt)) was invaluable to parsing out the data in a timely manner.

## Using the Prometheus exporter
To use the Prometheus exporter, simply download the binary and run it
```
./wx200
```

The following arguments are available

|Argument|Default|Description|
|-|-|-|
|`--listen-port`, `-p`|`9041`|Port that metrics server listens on. Metrics available at (`<host>:<listen-port>/metrics`)|
|`--com-port`, `-c`|`/dev/ttyUSB0`|COM port that the WX200 device is attached|
|`--verbosity`, `-v`|`info`|Log level (debug, info, warn)|

The exporter can also be run in Docker like so
```
docker run -d --rm -p 9041:9041 bartlettc/wx200 -c /dev/ttyUSB0 -p 9041
```

## Using the WX200 library
Basic example of using the WX200 golang library

```
import github.com/bartlettc22/wx200/pkg/wx200

...

    // Temperature data will be pushed to this channel
    temperatureDataChan := make(chan wx200.Temperature, 1)
	wx = wx200.New(&wx200.Config{
		ComPortName:         "/dev/ttyUSB0",
		TemperatureDataChan: temperatureDataChan,
	})

    // Starts async reading of serial data from the WX200
    go wx.Go()

    // Process incoming data
	for d := range temperatureDataChan {
        fmt.Printf("Indoor Temp: %d, Outdoor Temp: %d\n", d.Indoor, d.Outdoor)
	}

...
```

Channels are available for all sensors (see [main.go](main.go) for an more extensive example of how to use it).  Also, see the [godocs](https://godoc.org/github.com/bartlettc22/wx200/pkg/wx200) for more information on library.

## Notes
* Only some basic metrics are exposed via the Prometheus exporter but the library can read all of them in
* There seems to be a mistake in the serial mapping around dew point out of range values.  These are hard to reproduce so they may be wrong currently

[Build-Status-Url]: https://travis-ci.org/bartlettc22/wx200
[Build-Status-Image]: https://travis-ci.org/bartlettc22/wx200.svg?branch=master
[reportcard-url]: https://goreportcard.com/report/github.com/bartlettc22/wx200
[reportcard-image]: https://goreportcard.com/badge/github.com/bartlettc22/wx200
[godoc-url]: https://godoc.org/github.com/bartlettc22/wx200/pkg/wx200
[godoc-image]: https://godoc.org/github.com/bartlettc22/wx200pkg/wx200?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg