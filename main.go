package main

import "fmt"
import "log"
import "io"
import "github.com/jacobsa/go-serial/serial"
//import "github.com/davecgh/go-spew/spew"

func main() {
    // Set up options.
    options := serial.OpenOptions{
      PortName: "/dev/ttyUSB0",
      BaudRate: 9600,
      DataBits: 8,
      StopBits: 1,
      MinimumReadSize: 4,
    }

    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      log.Fatalf("serial.Open: %v", err)
    }

    // Make sure to close it later.
    defer port.Close()


    // Read in first byte to determine header

    b := make([]byte, 1)
    for {
	_, err := port.Read(b)
	//fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
//	fmt.Printf("b[:n] = %q\n", b[:n])
	if err == io.EOF {
		break
	}

	switch b[0] {
	case '\x8f':
		fmt.Println("Recieved Time/Humidity data...")
	case '\x9f':
		fmt.Println("Recieved Temperature data...")
	case '\xaf':
		fmt.Println("Recieved Barometer/Dew Point data...")
	case '\xbf':
		fmt.Println("Recieved Rain data...")
	case '\xcf':
		fmt.Println("Recieved Wind/Wind Chill/General data...")
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
    // Write 4 bytes to the port.
    // b := []byte{0x00, 0x01, 0x02, 0x03}
    //n, err := port.Write(b)
    //if err != nil {
    //  log.Fatalf("port.Write: %v", err)
    //}

    //fmt.Println("Wrote", n, "bytes.")
}
