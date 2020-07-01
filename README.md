# Radio Shack WX-200 Electronic Weather Station
[![Build][Build-Status-Image]][Build-Status-Url] [![ReportCard][reportcard-image]][reportcard-url] [![GoDoc][godoc-image]][godoc-url] [![License][license-image]][license-url]

This project contains the following:
* Prometheus exporter for WX200 readings (amd64/arm/Docker)
* Golang library for interfacing with the WX200 via serial communication


## Credits
This project could not have been completed without the work of Mike Wingstrom, Glynne Tolar and Tim Witham at [wx200.planetfall.com](http://wx200.planetfall.com/).  Their [serial protocol mapping](http://wx200.planetfall.com/wx200.txt) (also located [here](docs/wx200_serial_protocol.txt)) was invaluable to parsing out the data in a timely manner.

## Using the Prometheus exporter
* TBD

## Using the WX200 library
* TBD

## Notes
* There seems to be a mistake in the serial mapping around dew point out of range values.  These are hard to reproduce so they may be wrong currently

[Build-Status-Url]: https://travis-ci.org/bartlettc22/wx200
[Build-Status-Image]: https://travis-ci.org/bartlettc22/wx200.svg?branch=master
[reportcard-url]: https://goreportcard.com/report/github.com/bartlettc22/wx200
[reportcard-image]: https://goreportcard.com/badge/github.com/bartlettc22/wx200
[godoc-url]: https://godoc.org/github.com/bartlettc22/wx200/pkg/wx200
[godoc-image]: https://godoc.org/github.com/bartlettc22/wx200pkg/wx200?status.svg
[license-url]: http://opensource.org/licenses/MIT
[license-image]: https://img.shields.io/npm/l/express.svg