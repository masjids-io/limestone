package entity

type PrayerAdjustments struct {
	FajrAdjustment    int32 `gorm:"default:0"`
	DhuhrAdjustment   int32 `gorm:"default:0"`
	AsrAdjustment     int32 `gorm:"default:0"`
	MaghribAdjustment int32 `gorm:"default:0"`
	IshaAdjustment    int32 `gorm:"default:0"`
}

type CalculationMethod int64

const (
	OTHER CalculationMethod = iota
	MUSLIM_WORLD_LEAGUE
	EGYPTIAN
	KARACHI
	UMM_AL_QURA
	DUBAI
	MOON_SIGHTING_COMMITTEE
	NORTH_AMERICA
	KUWAIT
	QATAR
	SINGAPORE
	UOIF
)

type AsrJuristicMethod int64

const (
	SHAFI_HANBALI_MALIKI AsrJuristicMethod = iota
	HANAFI
)

type HighLatitudeRule int64

const (
	NO_HIGH_LATITUDE_RULE HighLatitudeRule = iota
	MIDDLE_OF_THE_NIGHT
	SEVENTH_OF_THE_NIGHT
	TWILIGHT_ANGLE
)

type PrayerTimesConfiguration struct {
	CalculationMethod CalculationMethod `sql:"type:ENUM('OTHER','MUSLIM_WORLD_LEAGUE','EGYPTIAN','KARACHI','UMM_AL_QURA','DUBAI','MOON_SIGHTING_COMMITTEE','NORTH_AMERICA','KUWAIT','QATAR','SINGAPORE','UOIF')" gorm:"column:calculation_method"`
	FajrAngle         float64           `gorm:"default:0"`
	IshaAngle         float64           `gorm:"default:0"`
	IshaInterval      int32             `gorm:"default:0"`
	AsrMethod         AsrJuristicMethod `sql:"type:ENUM('SHAFI_HANBALI_MALIKI','HANAFI')" gorm:"column:asr_method"`
	HighLatitudeRule  HighLatitudeRule  `sql:"type:ENUM('NO_HIGH_LATITUDE_RULE','MIDDLE_OF_THE_NIGHT','SEVENTH_OF_THE_NIGHT','TWILIGHT_ANGLE')" gorm:"column:high_latitude_rule"`
	Adjustments       PrayerAdjustments `gorm:"embedded"`
}
