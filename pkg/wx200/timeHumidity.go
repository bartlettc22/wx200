package wx200

import (
	"errors"
	"time"
)

type Time struct {
	Second       int8
	Minute       int8
	Hour         int8
	Day          int8
	HourFormat24 bool // Is time in 24-hour format
	DateFormatMD bool // Is date in month-day format (day-month, if false)
	Month        int8
	AlarmMinute  int8
	AlarmHour    int8
	AlarmSet     bool
}

// Humidity holds humidity values, records and alarm data
type Humidity struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Current indoor humidity (%)
	Indoor int8

	// Current indoor humidity out of range
	IndoorOR bool

	// Record high for indoor humidity (%)
	IndoorHi int8

	// Record high for indoor humidity out of range
	IndoorHiOR bool

	// Date of record high for indoor humidity
	IndoorHiDate time.Time

	// Threshhold for alarm for high indoor humidity (%)
	IndoorAlarmHi int8

	// Record low for indoor humidity (%)
	IndoorLo int8

	// Date of record low for indoor humidity
	IndoorLoDate time.Time

	// Threshhold for alerm for low indoor humidity (%)
	IndoorAlarmLo int8

	// Alarm for indoor humidity is set
	IndoorAlarmSet bool

	// Current outdoor humidity out of range
	OutdoorOR bool

	// Record high for outdoor humidity out of range
	OutdoorHiOR bool

	// Threshhold for alarm for high outdoor humidity (%)
	OutdoorAlarmHi int8

	// Threshhold for alerm for low outdoor humidity (%)
	OutdoorAlarmLo int8

	// Alarm for outdoor humidity is set
	OutdoorAlarmSet bool

	// Checksum
	ChecksumValid bool
}

func (w *WX200) ReadTimeHumidity() error {

	now := time.Now()
	err := readSerial(w.comPort, bufTimeHumidity)
	if err != nil {
		return err
	}

	t := Time{}
	t.Second = getCombinedDecimal(bufTimeHumidity[0])
	t.Minute = getCombinedDecimal(bufTimeHumidity[1])
	t.Hour = getCombinedDecimal(bufTimeHumidity[2])
	t.Day = getCombinedDecimal(bufTimeHumidity[3])
	s1, s2 := splitByte(bufTimeHumidity[4])
	t.HourFormat24 = isBitSet(s1, 0)
	t.DateFormatMD = isBitSet(s1, 1)
	t.Month = int8(s2)
	t.AlarmMinute = getCombinedDecimal(bufTimeHumidity[5])
	t.AlarmHour = getCombinedDecimal(bufTimeHumidity[6])

	var (
		indoorHiMonth  byte
		indoorLoMonth  byte
		indoorLo       [2]byte
		indoorHiMinute int8
		indoorLoMinute [2]byte
		indoorHiHour   int8
		indoorLoHour   [2]byte
		indoorHiDay    int8
		indoorLoDay    [2]byte
	)
	indoorHiMinute = getCombinedDecimal(bufTimeHumidity[9])
	indoorHiHour = getCombinedDecimal(bufTimeHumidity[10])
	indoorHiDay = getCombinedDecimal(bufTimeHumidity[11])
	indoorLo[1], indoorHiMonth = splitByte(bufTimeHumidity[12])
	indoorLoMinute[1], indoorLo[0] = splitByte(bufTimeHumidity[13])
	indoorLoHour[1], indoorLoMinute[0] = splitByte(bufTimeHumidity[14])
	indoorLoDay[1], indoorLoHour[0] = splitByte(bufTimeHumidity[15])
	indoorLoMonth, indoorLoDay[0] = splitByte(bufTimeHumidity[16])
	x, y := splitByte(bufTimeHumidity[31])
	x1, y1 := splitByte(bufTimeHumidity[32])

	h := Humidity{}
	h.LastDataRecieved = now
	h.Indoor = getCombinedDecimal(bufTimeHumidity[7])
	h.IndoorOR = isBitSet(x, 3)
	h.IndoorHi = getCombinedDecimal(bufTimeHumidity[8])
	h.IndoorHiOR = isBitSet(x, 2)
	h.IndoorLo = combineDecimal(indoorLo)
	h.IndoorLoDate = makeRecordDate(int(indoorLoMonth), int(combineDecimal(indoorLoDay)), int(combineDecimal(indoorLoHour)), int(combineDecimal(indoorLoMinute)))
	h.IndoorHiDate = makeRecordDate(int(indoorHiMonth), int(indoorHiDay), int(indoorHiHour), int(indoorHiMinute))
	h.IndoorAlarmHi = getCombinedDecimal(bufTimeHumidity[17])
	h.IndoorAlarmLo = getCombinedDecimal(bufTimeHumidity[18])
	h.IndoorAlarmSet = isBitSet(x1, 2) && isBitSet(x1, 3)

	h.OutdoorOR = isBitSet(x, 0)
	h.OutdoorHiOR = isBitSet(y, 3)
	h.OutdoorAlarmHi = getCombinedDecimal(bufTimeHumidity[29])
	h.OutdoorAlarmLo = getCombinedDecimal(bufTimeHumidity[30])
	h.OutdoorAlarmSet = isBitSet(x1, 0) && isBitSet(x1, 1)

	t.AlarmSet = isBitSet(y1, 3)
	// spew.Dump(t)
	// spew.Dump(h)

	w.Humidity = h

	// Validate checksum
	checksumValid := validateChecksum(HEADER_TIME_HUMIDITY, bufTimeHumidity)
	if !checksumValid {
		return errors.New("Checksum failed on group 'Time/Humidity')")
	}

	return nil
}

// Gets a 2-digit combined decimal value 0-99
func getCombinedDecimal(b byte) int8 {
	b1, b2 := splitByte(b)
	return combineDecimal([2]byte{b1, b2})
}

func combineDecimal(b [2]byte) int8 {
	return int8(b[0]*10 + b[1])
}

// Splits the byte up into two bytes
// byte[0] = first 4 bits
// byte[1] = last 4 bits
func splitByte(b byte) (byte, byte) {
	return b >> 4, b & '\x0f'
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
func validateChecksum(headerByte byte, data []byte) bool {

	// The expected checksum byte comes in as the last byte of data
	expectedChecksum := data[len(data)-1]

	// Add up the header + data (skipping the checksum byte)
	checkSum := int(headerByte)
	for _, val := range data[:len(data)-2] {
		checkSum = checkSum + int(val)
	}

	return int8(checkSum) == int8(expectedChecksum)

}
