package wx200

import (
	"time"
)

// General contains general WX200 display information
type General struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// PowerSourceDC
	// DCPower=true , ACPower=false
	PowerSourceDC bool

	// LowPowerIndicator
	// On=true , Off=false
	LowPowerIndicator bool

	// DisplaySelected indicates the display screen that is selected
	// Time=0, Temp=1, ... Rain=7
	DisplaySelected uint8

	// DisplaySubscreen indicates the display subscreen that is selected
	// First=0, Second=1, Third=2, Fourth=3
	DisplaySubscreen uint8

	// DisplayType indicates the display screen type that is selected
	// Main=0, Mem=1, AlarmIn=2, AlarmOut=3
	DisplayType uint8
}

const (
	DISPLAY_SELECTED_TIME = iota
	DISPLAY_SELECTED_TEMP
	DISPLAY_SELECTED_HUMIDITY
	DISPLAY_SELECTED_DEW
	DISPLAY_SELECTED_BARO
	DISPLAY_SELECTED_WIND
	DISPLAY_SELECTED_WINDCHILL
	DISPLAY_SELECTED_RAIN
)

const (
	DISPLAY_SUB_FIRST = iota
	DISPLAY_SUB_SECOND
	DISPLAY_SUB_THIRD
	DISPLAY_SUB_FOURTH
)

const (
	DISPLAY_TYPE_MAIN = iota
	DISPLAY_TYPE_MEM
	DISPLAY_TYPE_ALARM_IN
	DISPLAY_TYPE_ALARM_OUT
)

const (
	WIND_UNITS_MPH = iota
	WIND_UNITS_KNOTS
	WIND_UNITS_MS
	WIND_UNITS_KPH
)

// Wind contains wind speed/direction values, history and alarm information
type Wind struct {

	// LastDataRecieved contains the time the last data was received
	LastDataRecieved time.Time

	// Gust speed (0m/s to 56m/s)
	GustSpeed float64

	// Gust speed out of range
	GustSpeedOR bool

	// Gust direction (degrees)
	GustDirection uint16

	// Average wind speed (0m/s to 56m/s)
	AvgSpeed float64

	// Avgerage wind speed out of range
	AvgSpeedOR bool

	// Average wind speed direction (degrees)
	AvgDirection uint16

	// Record high wind speed (0m/s to 56m/s)
	HiSpeed float64

	// Record high wind speed out of range
	HiSpeedOR bool

	// Record high wind speed direction (degrees)
	HiDirection uint16

	// Record high wind speed date
	HiDate time.Time

	// High wind alarm threshold (0m/s to 56m/s)
	AlarmThreshold int8

	// High wind alarm set
	AlarmSet bool

	// Wind display units
	// mph=0, knots=1, m/s=2, kph=3
	DisplayUnits uint8
}

type WindChill struct {
	// Wind chill (-85C to 60C)
	Chill int8

	// Record low wind chill (-85C to 60C)
	ChillLo int8

	// Record low wind chill date
	ChillLoDate time.Time

	// Wind chill alarm threshold (-85C to 60C)
	ChillAlarmThreshold int16

	// Wind chill alarm set
	ChillAlarmSet bool
}

func (w *WX200) readWindGeneral() error {

	// now := time.Now()
	now := time.Now()
	buf, err := w.readSample(w.bufWindGeneral, header_wind_general)
	if err != nil {
		return err
	}

	// General
	general := General{}
	general.LastDataRecieved = now
	general.PowerSourceDC = isBitSet(buf[23][0], 2)
	general.LowPowerIndicator = isBitSet(buf[23][0], 3)
	general.DisplaySelected, err = SubDecimal(buf[24][0], 0, 2)
	w.error(err)
	general.DisplaySubscreen, err = SubDecimal(buf[24][1], 0, 1)
	w.error(err)
	general.DisplayType, err = SubDecimal(buf[24][1], 2, 3)
	w.error(err)
	if w.config.GeneralDataChan != nil {
		w.config.GeneralDataChan <- general
	}

	// Wind
	wind := Wind{}
	wind.LastDataRecieved = now
	wind.GustSpeed = float64(buf[2][1])*10 + (float64(combineDecimal(buf[1])) / 10)
	wind.GustSpeedOR = isBitSet(buf[25][0], 3)
	wind.GustDirection = uint16(combineDecimal(buf[3]))*10 + uint16(buf[2][0])
	wind.AvgSpeed = float64(buf[5][1])*10 + (float64(combineDecimal(buf[4])) / 10)
	wind.AvgSpeedOR = isBitSet(buf[25][0], 2)
	wind.AvgDirection = uint16(combineDecimal(buf[6]))*10 + uint16(buf[5][0])
	wind.HiSpeed = float64(buf[8][1])*10 + (float64(combineDecimal(buf[7])) / 10)
	wind.HiSpeedOR = isBitSet(buf[25][0], 1)
	wind.HiDirection = uint16(combineDecimal(buf[9]))*10 + uint16(buf[8][0])
	wind.HiDate = makeRecordDate(int(buf[13][1]), int(combineDecimal(buf[12])), int(combineDecimal(buf[11])), int(combineDecimal(buf[10])))
	wind.AlarmThreshold = int8(buf[14][1])*10 + int8(buf[13][0])
	wind.AlarmSet = isBitSet(buf[25][1], 2)
	wind.DisplayUnits, err = SubDecimal(buf[15][0], 2, 3)
	w.error(err)
	if w.config.WindDataChan != nil {
		w.config.WindDataChan <- wind
	}

	// Wind Chill
	chill := WindChill{}
	chillSign := int8(1)
	if isBitSet(buf[21][0], 1) {
		chillSign = -1
	}
	chill.Chill = int8(combineDecimal(buf[16])) * chillSign
	chillLoSign := int8(1)
	if isBitSet(buf[21][0], 0) {
		chillLoSign = -1
	}
	chill.ChillLo = int8(combineDecimal(buf[17])) * chillLoSign
	chill.ChillLoDate = makeRecordDate(int(buf[21][1]), int(combineDecimal(buf[20])), int(combineDecimal(buf[19])), int(combineDecimal(buf[18])))
	chillAlarmThresholdMultiplier, err := SubDecimal(buf[23][0], 0, 0)
	w.error(err)
	chillAlarmSign := int16(1)
	if isBitSet(buf[23][1], 3) {
		chillAlarmSign = -1
	}
	// Wind chill alarm data is in F for some reason so we convert to C
	chillAlarmF := (int16(chillAlarmThresholdMultiplier)*100 + int16(combineDecimal(buf[22]))) * chillAlarmSign
	chill.ChillAlarmThreshold = int16(float32(chillAlarmF-32.0) * (float32(5) / float32(9)))
	chill.ChillAlarmSet = isBitSet(buf[25][1], 1)
	if w.config.WindChillDataChan != nil {
		w.config.WindChillDataChan <- chill
	}

	return nil
}
