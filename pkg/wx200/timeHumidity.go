package wx200

import (
	"fmt"
	"time"
)

// Time display hour format
const (
	DISPLAY_HOUR_FORMAT_12 = iota
	DISPLAY_HOUR_FORMAT_24
)

// Time display date format
const (
	DISPLAY_DATE_FORMAT_DM = iota
	DISPLAY_DATE_FORMAT_MD
)

// Time holds time values, formats and alarm data
type Time struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Date/Time set on the WX200 device
	Date time.Time

	// Format of hour
	// 12hr=0 24hr=1
	HourFormat uint8

	// Format of date
	// day-month=0 month-day=1
	DateFormat uint8

	// Alarm time (the d/m/y part represent the next time the alarm will go off)
	Alarm time.Time

	// Alarm is set
	AlarmSet bool
}

// Humidity holds humidity values, records and alarm data
type Humidity struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Current indoor humidity (10% to 97%)
	Indoor uint8

	// Current indoor humidity out of range
	IndoorOR bool

	// Record high for indoor humidity (10% to 97%)
	IndoorHi uint8

	// Record high for indoor humidity out of range
	IndoorHiOR bool

	// Date of record high for indoor humidity
	IndoorHiDate time.Time

	// Threshold for alarm for high indoor humidity (10% to 97%)
	IndoorHiAlarm uint8

	// Record low for indoor humidity (10% to 97%)
	IndoorLo uint8

	// Date of record low for indoor humidity
	IndoorLoDate time.Time

	// Threshold for alerm for low indoor humidity (10% to 97%)
	IndoorLoAlarm uint8

	// Alarm for indoor humidity is set
	IndoorAlarmSet bool

	// Current outdoor humidity (10% to 97%)
	Outdoor uint8

	// Current outdoor humidity out of range
	OutdoorOR bool

	// Record high for outdoor humidity (10% to 97%)
	OutdoorHi uint8

	// Record high for outdoor humidity out of range
	OutdoorHiOR bool

	// Date of record high for outdoor humidity
	OutdoorHiDate time.Time

	// Threshold for alarm for high outdoor humidity 10% to 97%)
	OutdoorHiAlarm uint8

	// Record low for outdoor humidity (10% to 97%)
	OutdoorLo uint8

	// Date of record low for outdoor humidity
	OutdoorLoDate time.Time

	// Threshold for alerm for low outdoor humidity (10% to 97%)
	OutdoorLoAlarm uint8

	// Alarm for outdoor humidity is set
	OutdoorAlarmSet bool
}

func (w *WX200) readTimeHumidity() error {

	now := time.Now()
	buf, err := w.readSample(w.bufTimeHumidity, headerTimeHumidity)
	if err != nil {
		return err
	}

	t := Time{}
	t.LastDataRecieved = now
	// Since there is no "year", we assume it's the current year
	t.Date = time.Date(now.Year(), time.Month(buf[5][1]), int(combineDecimal(buf[4])), int(combineDecimal(buf[3])), int(combineDecimal(buf[2])), int(combineDecimal(buf[1])), 0, time.UTC)
	t.HourFormat, err = subDecimal(buf[5][0], 0, 0)
	w.error(err)
	t.DateFormat, err = subDecimal(buf[5][0], 1, 1)
	w.error(err)
	// For the alarm, we set the date as the next date the alarm will fire since only hour/minute are available on the device
	t.Alarm = time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), int(combineDecimal(buf[7])), int(combineDecimal(buf[6])), 0, 0, time.UTC)
	if t.Date.After(t.Alarm) {
		// Alarm already passed today, add a day
		t.Alarm = t.Alarm.Add(24 * time.Hour)
	}
	t.AlarmSet = isBitSet(buf[33][1], 3)
	if w.config.TimeDataChan != nil {
		select {
		case w.config.TimeDataChan <- t:
		default:
			return fmt.Errorf("Time data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	h := Humidity{}
	h.LastDataRecieved = now
	h.Indoor = combineDecimal(buf[8])
	h.IndoorOR = isBitSet(buf[32][0], 3)
	h.IndoorHi = combineDecimal(buf[9])
	h.IndoorHiOR = isBitSet(buf[32][0], 2)
	h.IndoorHiDate = makeRecordDate(int(buf[13][1]), int(combineDecimal(buf[12])), int(combineDecimal(buf[11])), int(combineDecimal(buf[10])))
	h.IndoorHiAlarm = combineDecimal(buf[18])
	h.IndoorLo = combineDecimal([2]byte{buf[14][1], buf[13][0]})
	h.IndoorLoDate = makeRecordDate(int(buf[17][0]), int(combineDecimal([2]byte{buf[17][1], buf[16][0]})), int(combineDecimal([2]byte{buf[16][1], buf[15][0]})), int(combineDecimal([2]byte{buf[15][1], buf[14][0]})))
	h.IndoorLoAlarm = combineDecimal(buf[19])
	h.IndoorAlarmSet = isBitSet(buf[33][0], 2) && isBitSet(buf[33][0], 3)
	h.Outdoor = combineDecimal(buf[20])
	h.OutdoorHi = combineDecimal(buf[21])
	h.OutdoorOR = isBitSet(buf[32][0], 0)
	h.OutdoorHiOR = isBitSet(buf[32][1], 3)
	h.OutdoorHiDate = makeRecordDate(int(buf[25][1]), int(combineDecimal(buf[24])), int(combineDecimal(buf[23])), int(combineDecimal(buf[22])))
	h.OutdoorHiAlarm = combineDecimal(buf[30])
	h.OutdoorLo = combineDecimal([2]byte{buf[26][1], buf[25][0]})
	h.OutdoorLoDate = makeRecordDate(int(buf[17][0]), int(combineDecimal([2]byte{buf[29][1], buf[28][0]})), int(combineDecimal([2]byte{buf[28][1], buf[27][0]})), int(int(combineDecimal([2]byte{buf[27][1], buf[26][0]}))))
	h.OutdoorLoAlarm = combineDecimal(buf[31])
	h.OutdoorAlarmSet = isBitSet(buf[33][0], 0) && isBitSet(buf[33][0], 1)

	if w.config.HumidityDataChan != nil {
		select {
		case w.config.HumidityDataChan <- h:
		default:
			return fmt.Errorf("Humidity data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	return nil
}
