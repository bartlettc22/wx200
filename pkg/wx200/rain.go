package wx200

import (
	// "errors"
	// "github.com/davecgh/go-spew/spew"
	"time"
)

// Rain contains rainfall values, history and alarm information
type Rain struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Rate of rainfall (0mm/hr to 998mm/hr)
	Rate uint16
}

func (w *WX200) readRain() error {

	// now := time.Now()
	// err := readSerial(w.comPort, w.bufRain)
	// if err != nil {
	// 	return err
	// }
	// w.info.SamplesRecieved = w.info.SamplesRecieved + 1

	// // Validate checksum
	// checksumValid := validateChecksum(HEADER_RAIN, w.bufRain)
	// if !checksumValid {
	// 	w.info.ChecksumFailures = w.info.ChecksumFailures + 1
	// 	return errors.New("Checksum failed on group 'Rain'")
	// }

	// // Split buffer up into 4-bit chunks to make it easier to work with
	// buf := chopBuffer(w.bufRain)

	now := time.Now()
	buf, err := w.readSample(w.bufRain, headerRain)
	if err != nil {
		return err
	}

	// Rain Values
	rain := Rain{}
	rain.LastDataRecieved = now
	rain.Rate = uint16(buf[2][1])*100 + uint16(combineDecimal(buf[1]))

	if w.config.RainDataChan != nil {
		w.config.RainDataChan <- rain
	}

	return nil
}
