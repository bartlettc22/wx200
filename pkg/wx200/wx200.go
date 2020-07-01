package wx200

import (
	"errors"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	// "time"
)

const (
	HEADER_TIME_HUMIDITY = 0x8f
	HEADER_TEMPERATURE   = 0x9f
	HEADER_BARO_DEW      = 0xaf
	HEADER_RAIN          = 0xbf
	HEADER_WIND_GENERAL  = 0xcf
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

	// Open the port.
	w.comPort, err = serial.Open(serialOptions)
	if err != nil {
		w.config.ErrorChan <- err
	}

	// Make sure to close it later.
	defer w.comPort.Close()

	headerByte := make([]byte, 1)

	// main loop for reading serial data
	// Time/Humidity 10 seconds
	// Temperature 10 seconds
	// Barometer/Dew Point 10 seconds
	// Rain 10 seconds
	// Wind/General 5 seconds
	for {

		// Read in first byte to determine header
		_, err := w.comPort.Read(headerByte)
		if err != nil {
			w.error(errors.New(fmt.Sprintf("Error reading serial data from %s", w.config.ComPortName)))
			continue
		}

		switch headerByte[0] {
		case HEADER_TIME_HUMIDITY:
			err = w.ReadTimeHumidity()
			// w.info.SamplesRecieved = w.info.SamplesRecieved + 1
		case HEADER_TEMPERATURE:
			err = w.readTemperature()
			// fmt.Println("Recieved Temperature data...")
			// w.info.SamplesRecieved = w.info.SamplesRecieved + 1
		case HEADER_BARO_DEW:
			err = w.readBaroDew()
			// fmt.Println("Recieved Barometer/Dew Point data...")
			// w.info.SamplesRecieved = w.info.SamplesRecieved + 1
		case HEADER_RAIN:
			err = w.readRain()
			// w.info.SamplesRecieved = w.info.SamplesRecieved + 1
		case HEADER_WIND_GENERAL:
			err = w.readWindGeneral()
			// w.info.SamplesRecieved = w.info.SamplesRecieved + 1
		default:
			// err = errors.New(fmt.Sprintf("Recieved unknown serial data header x%02x...", headerByte[0]))
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
			log.Fatalf("Fatal error occured and error handling not set: %v", err)
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
