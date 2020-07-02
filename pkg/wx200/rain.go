package wx200

import (
	"fmt"
	"time"
)

// Rainfall units enumerations
const (
	RAIN_UNITS_MM = iota
	RAIN_UNITS_INCHES
)

// Rain contains rainfall values, history and alarm information
type Rain struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Rate of rainfall (0mm/hr to 998mm/hr)
	Rate uint16

	// Rate of rainfall out of range
	RateOR bool

	// Yesterday's rainfall (0mm to 9999mm)
	Yesterday uint16

	// Total rainfall (0mm to 9999mm)
	Total uint16

	// Last reset
	Reset time.Time

	// Rainfall display units
	// mm=0, inches=1
	DisplayUnits uint8

	// High rain rate threshold (0mm/hr to 998mm/hr)
	AlarmThreshold uint16

	// High rain rate alarm set
	AlarmSet bool
}

func (w *WX200) readRain() error {

	now := time.Now()
	buf, err := w.readSample(w.bufRain, headerRain)
	if err != nil {
		return err
	}

	// Rain Values
	rain := Rain{}
	rain.LastDataRecieved = now
	rain.Rate = uint16(buf[2][1])*100 + uint16(combineDecimal(buf[1]))
	rain.RateOR = isBitSet(buf[12][0], 3)
	rain.Yesterday = uint16(combineDecimal(buf[4]))*100 + uint16(combineDecimal(buf[3]))
	rain.Total = uint16(combineDecimal(buf[6]))*100 + uint16(combineDecimal(buf[5]))
	rain.Reset = makeRecordDate(int(buf[10][1]), int(combineDecimal(buf[9])), int(combineDecimal(buf[8])), int(combineDecimal(buf[7])))
	rain.DisplayUnits, err = subDecimal(buf[10][0], 1, 1)
	w.error(err)
	// Rainfall alarm data is in in/hr for some reason so we convert to mm/hr
	rain.AlarmThreshold = uint16((float32(buf[12][1]*10) + float32(combineDecimal(buf[11]))/10) * float32(25.4))
	rain.AlarmSet = isBitSet(buf[12][0], 0)

	if w.config.RainDataChan != nil {
		select {
		case w.config.RainDataChan <- rain:
		default:
			return fmt.Errorf("Rain data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	return nil
}
