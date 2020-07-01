package wx200

import (
	// "errors"
	// "github.com/davecgh/go-spew/spew"
	"time"
)

// Time holds time values, formats and alarm data
type Time struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	Second       uint8
	Minute       uint8
	Hour         uint8
	Day          uint8
	HourFormat24 bool // Is time in 24-hour format
	DateFormatMD bool // Is date in month-day format (day-month, if false)
	Month        uint8
	AlarmMinute  uint8
	AlarmHour    uint8
	AlarmSet     bool
}

// Humidity holds humidity values, records and alarm data
type Humidity struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Current indoor humidity (%)
	Indoor uint8

	// Current indoor humidity out of range
	IndoorOR bool

	// Record high for indoor humidity (%)
	IndoorHi uint8

	// Record high for indoor humidity out of range
	IndoorHiOR bool

	// Date of record high for indoor humidity
	IndoorHiDate time.Time

	// Threshold for alarm for high indoor humidity (%)
	IndoorAlarmHi uint8

	// Record low for indoor humidity (%)
	IndoorLo uint8

	// Date of record low for indoor humidity
	IndoorLoDate time.Time

	// Threshold for alerm for low indoor humidity (%)
	IndoorAlarmLo uint8

	// Alarm for indoor humidity is set
	IndoorAlarmSet bool

	// Current outdoor humidity (%)
	Outdoor uint8

	// Current outdoor humidity out of range
	OutdoorOR bool

	// Record high for outdoor humidity (%)
	OutdoorHi uint8

	// Record high for outdoor humidity out of range
	OutdoorHiOR bool

	// Date of record high for outdoor humidity
	OutdoorHiDate time.Time

	// Threshold for alarm for high outdoor humidity (%)
	OutdoorAlarmHi uint8

	// Record low for outdoor humidity (%)
	OutdoorLo uint8

	// Date of record low for outdoor humidity
	OutdoorLoDate time.Time

	// Threshold for alerm for low outdoor humidity (%)
	OutdoorAlarmLo uint8

	// Alarm for outdoor humidity is set
	OutdoorAlarmSet bool
}

func (w *WX200) readTimeHumidity() error {

	now := time.Now()
	_, err := w.readSample(w.bufTimeHumidity, header_time_humidity)
	if err != nil {
		return err
	}
	// err := readSerial(w.comPort, w.bufTimeHumidity)
	// if err != nil {
	// 	return err
	// }

	// // Validate checksum
	// checksumValid := validateChecksum(HEADER_TIME_HUMIDITY, w.bufTimeHumidity)
	// if !checksumValid {
	// 	return errors.New("Checksum failed on group 'Time/Humidity'")
	// }

	t := Time{}
	t.LastDataRecieved = now
	t.Second = getCombinedDecimal(w.bufTimeHumidity[0])
	t.Minute = getCombinedDecimal(w.bufTimeHumidity[1])
	t.Hour = getCombinedDecimal(w.bufTimeHumidity[2])
	t.Day = getCombinedDecimal(w.bufTimeHumidity[3])
	s1, s2 := splitByte(w.bufTimeHumidity[4])
	t.HourFormat24 = isBitSet(s1, 0)
	t.DateFormatMD = isBitSet(s1, 1)
	t.Month = uint8(s2)
	t.AlarmMinute = getCombinedDecimal(w.bufTimeHumidity[5])
	t.AlarmHour = getCombinedDecimal(w.bufTimeHumidity[6])

	var (
		indoorHiMonth  byte
		indoorLoMonth  byte
		indoorLo       [2]byte
		indoorHiMinute uint8
		indoorLoMinute [2]byte
		indoorHiHour   uint8
		indoorLoHour   [2]byte
		indoorHiDay    uint8
		indoorLoDay    [2]byte
	)
	indoorHiMinute = getCombinedDecimal(w.bufTimeHumidity[9])
	indoorHiHour = getCombinedDecimal(w.bufTimeHumidity[10])
	indoorHiDay = getCombinedDecimal(w.bufTimeHumidity[11])
	indoorLo[1], indoorHiMonth = splitByte(w.bufTimeHumidity[12])
	indoorLoMinute[1], indoorLo[0] = splitByte(w.bufTimeHumidity[13])
	indoorLoHour[1], indoorLoMinute[0] = splitByte(w.bufTimeHumidity[14])
	indoorLoDay[1], indoorLoHour[0] = splitByte(w.bufTimeHumidity[15])
	indoorLoMonth, indoorLoDay[0] = splitByte(w.bufTimeHumidity[16])
	x, y := splitByte(w.bufTimeHumidity[31])
	x1, y1 := splitByte(w.bufTimeHumidity[32])

	h := Humidity{}
	h.LastDataRecieved = now
	h.Indoor = getCombinedDecimal(w.bufTimeHumidity[7])
	h.IndoorOR = isBitSet(x, 3)
	h.IndoorHi = getCombinedDecimal(w.bufTimeHumidity[8])
	h.IndoorHiOR = isBitSet(x, 2)
	h.IndoorLo = combineDecimal(indoorLo)
	h.IndoorLoDate = makeRecordDate(int(indoorLoMonth), int(combineDecimal(indoorLoDay)), int(combineDecimal(indoorLoHour)), int(combineDecimal(indoorLoMinute)))
	h.IndoorHiDate = makeRecordDate(int(indoorHiMonth), int(indoorHiDay), int(indoorHiHour), int(indoorHiMinute))
	h.IndoorAlarmHi = getCombinedDecimal(w.bufTimeHumidity[17])
	h.IndoorAlarmLo = getCombinedDecimal(w.bufTimeHumidity[18])
	h.IndoorAlarmSet = isBitSet(x1, 2) && isBitSet(x1, 3)
	h.Outdoor = getCombinedDecimal(w.bufTimeHumidity[19])
	h.OutdoorHi = getCombinedDecimal(w.bufTimeHumidity[20])
	h.OutdoorOR = isBitSet(x, 0)
	h.OutdoorHiOR = isBitSet(y, 3)
	outdoorHiMinute := getCombinedDecimal(w.bufTimeHumidity[21])
	outdoorHiHour := getCombinedDecimal(w.bufTimeHumidity[22])
	outdoorHiDay := getCombinedDecimal(w.bufTimeHumidity[23])
	ho0, outdoorHiMonth := splitByte(w.bufTimeHumidity[24])
	ho1, ho2 := splitByte(w.bufTimeHumidity[25])
	ho3, ho4 := splitByte(w.bufTimeHumidity[26])
	ho5, ho6 := splitByte(w.bufTimeHumidity[27])
	outdoorLoMonth, ho7 := splitByte(w.bufTimeHumidity[28])
	h.OutdoorLo = combineDecimal([2]byte{ho2, ho0})
	outdoorLoMinute := combineDecimal([2]byte{ho4, ho1})
	outdoorLoHour := combineDecimal([2]byte{ho6, ho3})
	outdoorLoDay := combineDecimal([2]byte{ho7, ho5})
	h.OutdoorLoDate = makeRecordDate(int(outdoorLoMonth), int(outdoorLoDay), int(outdoorLoHour), int(outdoorLoMinute))

	h.OutdoorHiDate = makeRecordDate(int(outdoorHiMonth), int(outdoorHiDay), int(outdoorHiHour), int(outdoorHiMinute))
	h.OutdoorAlarmHi = getCombinedDecimal(w.bufTimeHumidity[29])
	h.OutdoorAlarmLo = getCombinedDecimal(w.bufTimeHumidity[30])
	h.OutdoorAlarmSet = isBitSet(x1, 0) && isBitSet(x1, 1)

	t.AlarmSet = isBitSet(y1, 3)

	// Push our values to their respective data channels
	if w.config.TimeDataChan != nil {
		w.config.TimeDataChan <- t
	}
	if w.config.HumidityDataChan != nil {
		w.config.HumidityDataChan <- h
	}

	return nil
}
