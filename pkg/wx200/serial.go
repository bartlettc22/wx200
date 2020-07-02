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
