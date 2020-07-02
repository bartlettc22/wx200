package wx200

import (
	"fmt"
	"time"
)

// Barometer units enumerations
const (
	BARO_UNITS_IN = iota
	BARO_UNITS_MM
	BARO_UNITS_MB
	BARO_UNITS_HPA
)

// Barometer trend enumerations
const (
	BARO_TREND_RISING  = 1
	BARO_TREND_STEADY  = 2
	BARO_TREND_FALLING = 4
)

// Barometer prediction enumerations
const (
	BARO_PREDICTION_SUNNY  = 1
	BARO_PREDICTION_CLOUDY = 2
	BARO_PREDICTION_PARTLY = 4
	BARO_PREDICTION_RAINY  = 8
)

// Barometer contains pressure values, history and alarm information
type Barometer struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Local contains the local pressure (795mb to 1050mb)
	Local uint16

	// SeaLevel contains sea-level pressure (795mb to 1050mb)
	SeaLevel float32

	// Barometer display units
	// inches=0, mm=1, mb=2, hpa=3
	DisplayUnits uint8

	// Trend indicates the change in pressure
	// rising=1, steady=2, falling=4
	Trend uint8

	// Weather prediction
	// sunny=1, cloudy=2, partly=4, rainy=8
	Prediction uint8

	// Barometer alarm threshold (1mb to 16mb)
	AlarmThreshold uint8

	// Barometer alarm set
	AlarmSet bool
}

// DewPoint contains dew point values, history and alarm information
type DewPoint struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Indoor dew point (0C to 47C)
	Indoor uint8

	// Indoor dew point out of range
	IndoorOR bool

	// Record high indoor dew point (0C to 47C)
	IndoorHi uint8

	// Record high indoor dew point date
	IndoorHiDate time.Time

	// Record low indoor dew point (0C to 47C)
	IndoorLo uint8

	// Record low indoor dew point out of range
	IndoorLoOR bool

	// Record low indoor dew point date
	IndoorLoDate time.Time

	// Indoor dew point alarm threshold (1C to 16C)
	IndoorAlarmThreshold uint8

	// Outdoor dew point (0C to 56C)
	Outdoor uint8

	// Outdoor dew point out of range
	OutdoorOR bool

	// Record high outdoor dew point (0C to 47C)
	OutdoorHi uint8

	// Record high outdoor dew point date
	OutdoorHiDate time.Time

	// Record low outdoor dew point (0C to 47C)
	OutdoorLo uint8

	// Record low outdoor dew point out of range
	OutdoorLoOR bool

	// Record low outdoor dew point date
	OutdoorLoDate time.Time

	// Outdoor dew point alarm threshold (1C to 16C)
	OutdoorAlarmThreshold uint8

	// Dew point alarm is set
	AlarmSet bool
}

func (w *WX200) readBaroDew() error {

	now := time.Now()
	buf, err := w.readSample(w.bufBaroDew, headerBaroDew)
	if err != nil {
		return err
	}

	// Barometer
	baro := Barometer{}
	baro.LastDataRecieved = now
	baro.Local = uint16(combineDecimal(buf[2]))*100 + uint16(combineDecimal(buf[1]))
	baro.SeaLevel = float32(buf[5][1])*1000 + float32(combineDecimal(buf[4]))*10 + float32(combineDecimal(buf[3]))/10
	baro.DisplayUnits, err = subDecimal(buf[5][0], 0, 1)
	w.error(err)
	baro.Trend, err = subDecimal(buf[6][0], 0, 2)
	w.error(err)
	baro.Prediction = uint8(buf[6][1])
	// Threshold comes out indexed at 0 so we have to add one
	baro.AlarmThreshold = uint8(buf[29][1]) + 1
	baro.AlarmSet = isBitSet(buf[29][0], 3)
	if w.config.BarometerDataChan != nil {
		select {
		case w.config.BarometerDataChan <- baro:
		default:
			return fmt.Errorf("Barometer data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	// Dew Point
	dew := DewPoint{}
	dew.LastDataRecieved = now
	dew.Indoor = combineDecimal(buf[7])
	// serial proto mapping doc has error - I believe what is below is correct
	dew.IndoorOR = isBitSet(buf[28][0], 1)
	dew.IndoorHi = combineDecimal(buf[8])
	dew.IndoorHiDate = makeRecordDate(int(buf[12][1]), int(combineDecimal(buf[11])), int(combineDecimal(buf[10])), int(combineDecimal(buf[9])))
	dew.IndoorLo = combineDecimal([2]byte{buf[13][1], buf[12][0]})
	// serial proto mapping doc has error - I believe what is below is correct
	dew.IndoorLoOR = isBitSet(buf[28][1], 3)
	dew.IndoorLoDate = makeRecordDate(int(buf[16][0]), int(combineDecimal([2]byte{buf[16][1], buf[15][0]})), int(combineDecimal([2]byte{buf[15][1], buf[14][0]})), int(combineDecimal([2]byte{buf[14][1], buf[13][0]})))
	// Threshold comes out indexed at 0 so we have to add one
	dew.IndoorAlarmThreshold = uint8(buf[17][1]) + 1

	dew.Outdoor = combineDecimal(buf[18])
	// serial proto mapping doc has error - I believe what is below is correct
	dew.OutdoorOR = isBitSet(buf[28][1], 2)
	dew.OutdoorHi = combineDecimal(buf[19])
	dew.OutdoorHiDate = makeRecordDate(int(buf[23][1]), int(combineDecimal(buf[22])), int(combineDecimal(buf[21])), int(combineDecimal(buf[20])))
	dew.OutdoorLo = combineDecimal([2]byte{buf[24][1], buf[23][0]})
	// serial proto mapping doc has error - I believe what is below is correct
	dew.OutdoorLoOR = isBitSet(buf[28][1], 0)
	dew.OutdoorLoDate = makeRecordDate(int(buf[27][0]), int(combineDecimal([2]byte{buf[27][1], buf[26][0]})), int(combineDecimal([2]byte{buf[26][1], buf[25][0]})), int(combineDecimal([2]byte{buf[25][1], buf[24][0]})))
	// Threshold comes out indexed at 0 so we have to add one
	dew.OutdoorAlarmThreshold = uint8(buf[17][0]) + 1
	dew.AlarmSet = isBitSet(buf[29][0], 1) && isBitSet(buf[29][0], 2)
	if w.config.DewPointDataChan != nil {
		select {
		case w.config.DewPointDataChan <- dew:
		default:
			return fmt.Errorf("Dew Point data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	return nil
}
