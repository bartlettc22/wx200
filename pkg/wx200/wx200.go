package wx200

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"time"
)

const (
	headerTimeHumidity = 0x8f
	headerTemperature  = 0x9f
	headerBaroDew      = 0xaf
	headerRain         = 0xbf
	headerWindGeneral  = 0xcf
)

// Config contains configuration parameters for the WX200 library
type Config struct {
	ComPortName         string
	BarometerDataChan   chan Barometer
	InfoDataChan        chan Info
	TimeDataChan        chan Time
	HumidityDataChan    chan Humidity
	WindDataChan        chan Wind
	WindChillDataChan   chan WindChill
	GeneralDataChan     chan General
	RainDataChan        chan Rain
	TemperatureDataChan chan Temperature
	DewPointDataChan    chan DewPoint
	ErrorChan           chan error
}

// WX200 is the main library class
type WX200 struct {
	config          *Config
	comPort         io.ReadWriteCloser
	bufTimeHumidity []byte
	bufWindGeneral  []byte
	bufRain         []byte
	bufTemperature  []byte
	bufBaroDew      []byte
	info            Info
	ChecksumErrors  int64
}

// New constructs a new WX200 based on the provided Config
func New(config *Config) *WX200 {
	return &WX200{
		config:          config,
		bufTimeHumidity: make([]byte, 34),
		bufWindGeneral:  make([]byte, 26),
		bufRain:         make([]byte, 13),
		bufTemperature:  make([]byte, 33),
		bufBaroDew:      make([]byte, 30),
		info:            Info{},
	}
}

// Go begins serial communication and reading of sample data
// Is meant to be run as a goroutine and passes samples back
// via any provided Config data channel values
// Time/Humidity is received every 10 seconds
// Temperature is received every 10 seconds
// Barometer/Dew Point is received every 10 seconds
// Rain is received every 10 seconds
// Wind/Wind Chill/General is received every 5 seconds
func (w *WX200) Go() {

	var err error

	serialOptions := serial.OpenOptions{
		PortName:        w.config.ComPortName,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	reconnect := false
	headerByte := make([]byte, 1)

	// main loop for reading serial data
	for {

		// Connect to serial port if not already connected
		if w.comPort == nil || reconnect == true {
			w.comPort, err = serial.Open(serialOptions)
			if err != nil {
				w.error(fmt.Errorf("Error opening serial communication to WX200: %v. Retrying...", err))
				time.Sleep(5 * time.Second)
				continue
			} else {
				defer w.comPort.Close()
				reconnect = false
			}
		}

		// Read in first byte to determine header
		_, err := w.comPort.Read(headerByte)
		if err != nil {
			w.error(fmt.Errorf("Error reading serial data: %v. Reconnecting...", err))
			reconnect = true
			continue
		}

		switch headerByte[0] {
		case headerTimeHumidity:
			err = w.readTimeHumidity()
		case headerTemperature:
			err = w.readTemperature()
		case headerBaroDew:
			err = w.readBaroDew()
		case headerRain:
			err = w.readRain()
		case headerWindGeneral:
			err = w.readWindGeneral()
		default:
			// Seems to be some unknown data that comes through, just skip and move on...
			err = nil
		}

		if w.config.InfoDataChan != nil {
			w.config.InfoDataChan <- w.info
		}

		// This should check and make sure its just checksum errors
		if err != nil {
			w.ChecksumErrors = w.ChecksumErrors + 1
		}
		w.error(err)

	}
}

// Ready returns true once samples have been read from all groups
// Since there are 5 sample groups (4 every 10 seconds and 1 every 5)
// We wait 2 sample cycles (up to 20 seconds) to make sure all samples are
// read and the channels have enough time to process
func (w *WX200) Ready() bool {
	if w.info.SamplesRecieved >= 12 {
		return true
	}
	return false
}

func (w *WX200) error(err error) {
	if err != nil {
		if w.config.ErrorChan != nil {
			w.config.ErrorChan <- err
		} else {
			log.Fatalf("Fatal error occurred and error handling not set: %v", err)
		}
	}
}
