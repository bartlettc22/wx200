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

type Config struct {
	ComPortName         string
	TimeDataChan        chan Time
	HumidityDataChan    chan Humidity
	WindDataChan        chan Wind
	GeneralDataChan     chan General
	WindChillDataChan   chan WindChill
	RainDataChan        chan Rain
	TemperatureDataChan chan Temperature
	BarometerDataChan   chan Barometer
	DewPointDataChan    chan DewPoint
	InfoDataChan        chan Info
	ErrorChan           chan error
}

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
	// Time/Humidity 10 seconds
	// Temperature 10 seconds
	// Barometer/Dew Point 10 seconds
	// Rain 10 seconds
	// Wind/General 5 seconds
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
			err = nil
			// err = errors.New(fmt.Sprintf("Received unknown serial data header x%02x...", headerByte[0]))
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
// We wait 2 sample cycles (up to 20 seconds)
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

// readSample reads in serial data equal to len(buf)
// returns a 2-D slice of the sample data
func (w *WX200) readSample(buf []byte, header byte) ([][2]byte, error) {

	err := readSerial(w.comPort, buf)
	if err != nil {
		return nil, err
	}
	w.info.SamplesRecieved = w.info.SamplesRecieved + 1

	// The expected checksum byte comes in as the last byte of data
	expectedChecksum := buf[len(buf)-1]

	// Add up the header + data (skipping the checksum byte)
	checkSum := uint(header)
	for _, val := range buf[:len(buf)-1] {
		checkSum = checkSum + uint(val)
	}

	if int8(checkSum) != int8(expectedChecksum) {
		w.info.ChecksumFailures = w.info.ChecksumFailures + 1
		return nil, fmt.Errorf("Checksum failed on group 0x%02d", header)
	}

	// Split buffer up into 4-bit chunks to make it easier to work with
	return chopBuffer(buf), nil

}
