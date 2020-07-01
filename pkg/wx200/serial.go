package wx200

import (
	"fmt"
	"io"
)

// readSerial fills up the provided buf with data from the reader
func readSerial(r io.ReadWriteCloser, buf []byte) error {
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return fmt.Errorf("Error reading %d bytes from serial: %v", len(buf), err)
	}
	// fmt.Printf("%d bytes read from serial\n", bytesRead)

	return nil
}
