package wx200

import (
	"fmt"
	"time"
)

// Temperature contains temperature values, history and alarm information
type Temperature struct {
	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time
}

func (w *WX200) readTemperature() error {

	now := time.Now()
	_, err := w.readSample(w.bufTemperature, headerTemperature)
	if err != nil {
		return err
	}

	temp := Temperature{}
	temp.LastDataRecieved = now
	if w.config.TemperatureDataChan != nil {
		select {
		case w.config.TemperatureDataChan <- temp:
		default:
			return fmt.Errorf("Temperature data cannot be sent to channel (might be full). Skipping sample.")
		}
	}

	return nil
}
