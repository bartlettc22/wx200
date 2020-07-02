package wx200

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math"
	"time"
)

// Temperature contains temperature values, history and alarm information
type Temperature struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Indoor temperature (-50C to 50C)
	Indoor float32

	// Record high indoor temperature (-50C to 50C)
	IndoorHi float32

	// Record high indoor temperature date
	IndoorHiDate time.Time

	// Indoor high temperature alarm (0C to 50C)
	IndoorHiAlarmThreshold uint8
}

func (w *WX200) readTemperature() error {

	now := time.Now()
	buf, err := w.readSample(w.bufTemperature, headerTemperature)
	if err != nil {
		return err
	}

	temp := Temperature{}
	temp.LastDataRecieved = now
	indoorSign := float32(1)
	if isBitSet(buf[2][1], 3) {
		indoorSign = -1
	}
	indoorMultiplier, err := subDecimal(buf[2][1], 0, 2)
	w.error(err)
	temp.Indoor = (float32(indoorMultiplier*10) + float32(combineDecimal(buf[1]))/float32(10)) * indoorSign
	indoorHiSign := float32(1)
	if isBitSet(buf[3][0], 3) {
		indoorHiSign = -1
	}
	indoorHiMultiplier, err := subDecimal(combineDecimal(buf[3]), 0, 6)
	w.error(err)
	temp.IndoorHi = (float32(indoorHiMultiplier) + float32(buf[2][0])/float32(10)) * indoorHiSign
	temp.IndoorHiDate = makeRecordDate(int(buf[7][1]), int(combineDecimal(buf[6])), int(combineDecimal(buf[5])), int(combineDecimal(buf[4])))
	// Temperature alarm data is in F for some reason so we convert to C
	indoorHiAlarmMultiplier, err := subDecimal(combineDecimal(buf[13]), 0, 4)
	w.error(err)
	indoorHiAlarmF := (int16(indoorHiAlarmMultiplier)*10 + int16(buf[12][0]))
	temp.IndoorHiAlarmThreshold = uint8(math.Round(float64(indoorHiAlarmF-32.0) * (float64(5) / float64(9))))

	spew.Dump(temp)

	if w.config.TemperatureDataChan != nil {
		select {
		case w.config.TemperatureDataChan <- temp:
		default:
			return fmt.Errorf("Temperature data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	return nil
}
