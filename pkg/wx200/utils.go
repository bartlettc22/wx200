package wx200

import (
	"errors"
	// "fmt"
	"time"
)

// Gets a 2-digit combined decimal value 0-99
// func getCombinedDecimal(b byte) uint8 {
// 	b1, b2 := splitByte(b)
// 	return combineDecimal([2]byte{b1, b2})
// }

func combineDecimal(b [2]byte) uint8 {
	return uint8(b[0]*10 + b[1])
}

// SubDecimal takes a byte b and extracts a subset of bits and returns the decimal number they represent
// For example subDecimal('\xf0', 1, 5) would take byte '\xf0' (represented as binary 11110000) and return
// the decimal value of positions 1-5 (11000) = 24
// startBit is the index of the starting bit (rightmost, starting at 0)
// endBit is the index of the ending bit (leftmost, ending at 7)
func SubDecimal(b byte, startBit uint8, endBit uint8) (uint8, error) {

	if endBit >= 8 {
		return 0, errors.New("Out of range: endBit must be less than 8")
	}

	if startBit > endBit {
		return 0, errors.New("Out of range: startBit must be less than or equal to endBit")
	}

	bitLength := endBit - startBit + 1
	mask := byte((1 << bitLength) - 1)
	return uint8((b >> startBit) & mask), nil
}

// Splits the byte up into two bytes
// byte[0] = first 4 bits
// byte[1] = last 4 bits
func splitByte(b byte) (byte, byte) {
	return b >> 4, b & '\x0f'
}

func splitByte2(b byte) [2]byte {
	return [2]byte{b >> 4, b & '\x0f'}
}

// combineByte takes two bytes (presumably with only 4 bits of data) and combines them to one
// 00000011, 00001100 becomes 00111100, etc.
func combineByte(b [2]byte) byte {
	return b[0]<<4 | b[1]
}

// isBitSet return true if the bit at position p is set
// p=0 is rightmost bit, p=7 is leftmost bit
func isBitSet(b byte, p int8) bool {

	if p >= 0 && p <= 7 {
		return (b>>p)&1 == 1
	}

	// out of bounds
	return false

}

// Presumes that if month/day has not already passed this year, that the record was set last year
func makeRecordDate(month int, day int, hour int, minute int) time.Time {

	now := time.Now()
	y := now.Year()
	if (time.Month(month) > now.Month()) || (time.Month(month) == now.Month() && day > now.Day()) {
		y = y - 1
	}

	return time.Date(y, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}

// validateChecksum validates that the data was passed over the serial line correctly
// This is done by adding up each byte of the group (including the group header) and
// comparing it to the checksum byte
// func validateChecksum(headerByte byte, data []byte) bool {

// 	// The expected checksum byte comes in as the last byte of data
// 	expectedChecksum := data[len(data)-1]
// 	// if headerByte == byte(HEADER_RAIN) {
// 	// 	fmt.Printf("Expected: hex:%02x int:%d\n", expectedChecksum, uint(expectedChecksum))
// 	// 	fmt.Printf("Adding %d bytes\n", uint(headerByte))
// 	// }

// 	// Add up the header + data (skipping the checksum byte)
// 	checkSum := uint(headerByte)
// 	// i := 1
// 	for _, val := range data[:len(data)-1] {
// 		checkSum = checkSum + uint(val)
// 		// i = i + 1
// 		// if headerByte == byte(HEADER_RAIN) {
// 		// 	fmt.Printf("Adding %d bytes, new total: %d\n", uint(val), checkSum)
// 		// }
// 	}
// 	// if headerByte == byte(HEADER_RAIN) {
// 	// 	fmt.Printf("Actual: int:%d hex:%02x uint8:%d\n", checkSum, checkSum, uint8(checkSum))
// 	// 	fmt.Printf("Checked %d bytes\n", i)
// 	// 	fmt.Println()

// 	// }
// 	return int8(checkSum) == int8(expectedChecksum)

// }

// Split into 4-bit chunks
// Pads output by adding zeros to first byte pair and shifts indexes up by 1
// so that they line up with docs
func chopBuffer(buf []byte) [][2]byte {
	out := make([][2]byte, len(buf)+1)

	out[0] = [2]byte{'\x00', '\x00'}

	for i, b := range buf {
		out[i+1] = splitByte2(b)
	}

	return out
}
