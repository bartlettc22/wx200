package wx200

import (
	"fmt"
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

	// Record low indoor temperature (-50C to 50C)
	IndoorLo float32

	// Record low indoor temperature date
	IndoorLoDate time.Time

	// Indoor low temperature alarm (0C to 50C)
	IndoorLoAlarmThreshold uint8

	// Indoor temperature alarm set
	IndoorAlarmSet bool

	// Outdoor temperature (-40C to 60C)
	Outdoor float32

	// Record high outdoor temperature (-40C to 60C)
	OutdoorHi float32

	// Record high outdoor temperature date
	OutdoorHiDate time.Time

	// Outdoor high temperature alarm (-40C to 60C)
	OutdoorHiAlarmThreshold int8

	// Record low outdoor temperature (-40C to 60C)
	OutdoorLo float32

	// Record low outdoor temperature date
	OutdoorLoDate time.Time

	// Outdoor low temperature alarm (-40C to 60C)
	OutdoorLoAlarmThreshold int8

	// Outdoor temperature alarm set
	OutdoorAlarmSet bool
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
	indoorHiMultiplier, err := subDecimal(combineDecimal(buf[3]), 0, 7)
	w.error(err)
	temp.IndoorHi = (float32(indoorHiMultiplier) + float32(buf[2][0])/float32(10)) * indoorHiSign
	temp.IndoorHiDate = makeRecordDate(int(buf[7][1]), int(combineDecimal(buf[6])), int(combineDecimal(buf[5])), int(combineDecimal(buf[4])))
	// Temperature alarm data is in F for some reason so we convert to C
	indoorHiAlarmMultiplier, err := subDecimal(combineDecimal(buf[13]), 0, 4)
	w.error(err)
	indoorHiAlarmF := (int16(indoorHiAlarmMultiplier)*10 + int16(buf[12][0]))
	temp.IndoorHiAlarmThreshold = uint8(math.Round(float64(indoorHiAlarmF-32.0) * (float64(5) / float64(9))))
	indoorLoSign := float32(1)
	if isBitSet(buf[8][0], 3) {
		indoorLoSign = -1
	}
	indoorLoMultiplier, err := subDecimal(combineDecimal(buf[8]), 0, 7)
	w.error(err)
	temp.IndoorLo = (float32(indoorLoMultiplier) + float32(buf[7][0])/float32(10)) * indoorLoSign
	temp.IndoorLoDate = makeRecordDate(int(buf[12][1]), int(combineDecimal(buf[11])), int(combineDecimal(buf[10])), int(combineDecimal(buf[9])))
	// Temperature alarm data is in F for some reason so we convert to C
	indoorLoAlarmMultiplier, err := subDecimal(buf[15][1], 0, 0)
	w.error(err)
	indoorLoAlarmF := indoorLoAlarmMultiplier*100 + combineDecimal(buf[14])
	temp.IndoorLoAlarmThreshold = uint8(math.Round(float64(indoorLoAlarmF-32.0) * (float64(5) / float64(9))))
	temp.IndoorAlarmSet = isBitSet(buf[32][0], 2) && isBitSet(buf[32][0], 3)

	outdoorSign := float32(1)
	if isBitSet(buf[17][1], 3) {
		outdoorSign = -1
	}
	outdoorMultiplier, err := subDecimal(buf[17][1], 0, 2)
	w.error(err)
	temp.Outdoor = (float32(outdoorMultiplier*10) + float32(combineDecimal(buf[16]))/float32(10)) * outdoorSign
	outdoorHiSign := float32(1)
	if isBitSet(buf[18][0], 3) {
		outdoorHiSign = -1
	}
	outdoorHiMultiplier, err := subDecimal(combineDecimal(buf[18]), 0, 7)
	w.error(err)
	temp.OutdoorHi = (float32(outdoorHiMultiplier) + float32(buf[17][0])/float32(10)) * outdoorHiSign
	temp.OutdoorHiDate = makeRecordDate(int(buf[22][1]), int(combineDecimal(buf[21])), int(combineDecimal(buf[20])), int(combineDecimal(buf[19])))
	// // Temperature alarm data is in F for some reason so we convert to C
	outdoorHiAlarmSign := int16(1)
	if isBitSet(buf[28][0], 3) {
		outdoorHiAlarmSign = -1
	}
	// Differs from protocol doc (0-3 vs 0-4)
	outdoorHiAlarmMultiplier, err := subDecimal(combineDecimal(buf[28]), 0, 3)
	w.error(err)
	outdoorHiAlarmF := (int16(outdoorHiAlarmMultiplier)*10 + int16(buf[27][0])) * outdoorHiAlarmSign
	// Temperature alarm data is in F for some reason so we convert to C
	temp.OutdoorHiAlarmThreshold = tempFToC(outdoorHiAlarmF)
	outdoorLoSign := float32(1)
	if isBitSet(buf[32][0], 3) {
		outdoorLoSign = -1
	}
	outdoorLoMultiplier, err := subDecimal(combineDecimal(buf[23]), 0, 7)
	w.error(err)
	temp.OutdoorLo = (float32(outdoorLoMultiplier) + float32(buf[22][0])/float32(10)) * outdoorLoSign
	temp.OutdoorLoDate = makeRecordDate(int(buf[27][1]), int(combineDecimal(buf[26])), int(combineDecimal(buf[25])), int(combineDecimal(buf[24])))
	outdoorLoAlarmSign := int16(1)
	if isBitSet(buf[30][1], 3) {
		outdoorLoAlarmSign = -1
	}
	outdoorLoAlarmMultiplier, err := subDecimal(buf[30][1], 0, 0)
	w.error(err)
	outdoorLoAlarmF := (int16(outdoorLoAlarmMultiplier)*100 + int16(combineDecimal(buf[29]))) * outdoorLoAlarmSign
	// Temperature alarm data is in F for some reason so we convert to C
	temp.OutdoorLoAlarmThreshold = tempFToC(outdoorLoAlarmF)
	temp.OutdoorAlarmSet = isBitSet(buf[32][0], 0) && isBitSet(buf[32][0], 1)

	if w.config.TemperatureDataChan != nil {
		select {
		case w.config.TemperatureDataChan <- temp:
		default:
			return fmt.Errorf("Temperature data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	return nil
}
