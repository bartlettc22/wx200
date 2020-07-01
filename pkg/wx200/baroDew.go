package wx200

import (
	"time"
)

type Barometer struct {
	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time
}

type DewPoint struct {
	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time
}

func (w *WX200) readBaroDew() error {

	now := time.Now()
	_, err := w.readSample(w.bufBaroDew, HEADER_BARO_DEW)
	if err != nil {
		return err
	}

	// Barometer
	baro := Barometer{}
	baro.LastDataRecieved = now

	// Dew Point
	dew := DewPoint{}
	dew.LastDataRecieved = now

	return nil
}
