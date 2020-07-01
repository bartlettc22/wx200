package wx200

import (
	"time"
)

type Temperature struct {
	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time
}

func (w *WX200) readTemperature() error {

	now := time.Now()
	_, err := w.readSample(w.bufTemperature, HEADER_TEMPERATURE)
	if err != nil {
		return err
	}

	temp := Temperature{}
	temp.LastDataRecieved = now

	return nil
}
