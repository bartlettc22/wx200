package wx200

import (
	// "fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	// "sync"
)

const (
	HEADER_TIME_HUMIDITY = '\x8f'
	HEADER_TEMPERATURE   = '\x9f'
	HEADER_BARO_DEW      = '\xaf'
	HEADER_RAIN          = '\xbf'
	HEADER_WIND_GENERAL  = '\xcf'
)

var bufTimeHumidity []byte

type Config struct {
	ComPortName string
	// StartupWG   sync.WaitGroup
}

type WX200 struct {
	config   *Config
	comPort  io.ReadWriteCloser
	Humidity Humidity
}

func New(config *Config) *WX200 {
	return &WX200{
		config:   config,
		Humidity: Humidity{},
	}
}

func (w *WX200) Go() {
	bufTimeHumidity = make([]byte, 34)
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
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer w.comPort.Close()

	// Read in first byte to determine header

	b := make([]byte, 1)
	// main loop for reading serial data
	// Time/Humidity 10 seconds
	// Temperature 10 seconds
	// Barometer/Dew Point 10 seconds
	// Rain 10 seconds
	// Wind/General 5 seconds
	for {
		_, err := w.comPort.Read(b)
		if err != nil {
			// Should send error to an error channel
		}
		//fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		//	fmt.Printf("b[:n] = %q\n", b[:n])
		// if err == io.EOF {
		// 	break
		// }

		switch b[0] {
		case HEADER_TIME_HUMIDITY:
			// fmt.Println("Recieved Time/Humidity data...")
			w.ReadTimeHumidity()
		case HEADER_TEMPERATURE:
			// fmt.Println("Recieved Temperature data...")
		case HEADER_BARO_DEW:
			// fmt.Println("Recieved Barometer/Dew Point data...")
		case HEADER_RAIN:
			// fmt.Println("Recieved Rain data...")
		case HEADER_WIND_GENERAL:
			// fmt.Println("Recieved Wind/General data...")
		default:
			//fmt.Println("Recieved Unknown data...")
		}

		// c := make([]byte, 34)
		// m, err := port.Read(c)
		// if err != nil {
		//   fmt.Errorf("%v", err)
		// }
		// spew.Dump(c[:m])
		//}

	}
}

// func New()
