package storage

// The manual adjustments to apply to the prayer timings. The value that each field is set to is
// the offset which will be added to the calculated time to obtain the final prayer time.
type PrayerAdjustments struct {
	// Adjustment offset for Fajr in minutes. Value can be negative.
	FajrAdjustment int32 `gorm:"default:0"`
	// Adjustment offset for Dhuhr in minutes. Value can be negative.
	DhuhrAdjustment int32 `gorm:"default:0"`
	// Adjustment offset for Asr in minutes. Value can be negative.
	AsrAdjustment int32 `gorm:"default:0"`
	// Adjustment offset for Maghrib in minutes. Value can be negative.
	MaghribAdjustment int32 `gorm:"default:0"`
	// Adjustment offset for Isha in minutes. Value can be negative.
	IshaAdjustment int32 `gorm:"default:0"`
}

// Defines the calculation method to use. If this field is set (excluding OTHER), then the Fajr
// and Isha angle fields are automatically set.
type CalculationMethod int64

const (
	OTHER CalculationMethod = iota
	// Muslim World League
	// Uses Fajr angle of 18 and an Isha angle of 17
	MUSLIM_WORLD_LEAGUE
	// Egyptian General Authority of Survey
	// Uses Fajr angle of 19.5 and an Isha angle of 17.5
	EGYPTIAN
	// University of Islamic Sciences, Karachi
	// Uses Fajr angle of 18 and an Isha angle of 18
	KARACHI
	// Umm al-Qura University, Makkah
	// Uses a Fajr angle of 18.5 and an Isha angle of 90. Note: You should add a +30 minute custom
	// adjustment of Isha during Ramadan.
	UMM_AL_QURA
	// The Gulf Region
	// Uses Fajr and Isha angles of 18.2 degrees.
	DUBAI
	// Moonsighting Committee
	// Uses a Fajr angle of 18 and an Isha angle of 18. Also uses seasonal adjustment values.
	MOON_SIGHTING_COMMITTEE
	// Referred to as the ISNA method
	// This method is included for completeness, but is not recommended.
	// Uses a Fajr angle of 15 and an Isha angle of 15.
	NORTH_AMERICA
	// Kuwait
	// Uses a Fajr angle of 18 and an Isha angle of 17.5
	KUWAIT
	// Qatar
	// Modified version of Umm al-Qura that uses a Fajr angle of 18.
	QATAR
	// Singapore
	// Uses a Fajr angle of 20 and an Isha angle of 18
	SINGAPORE
	// UOIF
	// Uses a Fajr angle of 12 and an Isha angle of 12
	UOIF
)

// The Juristic method to use for calculating Asr prayer times.
type AsrJuristicMethod int64

const (
	// Use the Shafi/Hanbali/Maliki method to calculate Asr timings.
	SHAFI_HANBALI_MALIKI AsrJuristicMethod = iota
	// Use the Hanafi method to calculate Asr timings.
	HANAFI
)

// The high latitude rule for calculating Fajr and Isha prayers.
type HighLatitudeRule int64

const (
	NO_HIGH_LATITUDE_RULE HighLatitudeRule = iota

	// Fajr will never be earlier than the middle of the night, and Isha will never be later than
	// the middle of the night.
	MIDDLE_OF_THE_NIGHT

	// Fajr will never be earlier than the beginning of the last seventh of the night, and Isha will
	// never be later than the end of the first seventh of the night.
	SEVENTH_OF_THE_NIGHT

	// Similar to SEVENTH_OF_THE_NIGHT, but instead of 1/7th, the fraction of the night used
	// is fajrAngle / 60 and ishaAngle/60.
	TWILIGHT_ANGLE
)

// A message that holds prayer times configuration. This message contains all the fields
// necessary to calculate prayer times.
type PrayerTimesConfiguration struct {
	// The calculation method to use.
	CalculationMethod CalculationMethod `sql:"type:ENUM('OTHER','MUSLIM_WORLD_LEAGUE','EGYPTIAN','KARACHI','UMM_AL_QURA','DUBAI','MOON_SIGHTING_COMMITTEE','NORTH_AMERICA','KUWAIT','QATAR','SINGAPORE','UOIF')" gorm:"column:calculation_method"`
	// The Fajr angle to use. This does not need to be set if the calculation method is set.
	// If both fields are set, then the calculation method field takes precedence.
	FajrAngle float64 `gorm:"default:0"`
	// The Isha angle to use. This does not need to be set if the calculation method is set.
	// If both fields are set, then the calculation method field takes precedence.
	IshaAngle float64 `gorm:"default:0"`
	// Minutes after Maghrib (if set, the time for Isha will be Maghrib plus the Isha interval).
	IshaInterval int32 `gorm:"default:0"`
	// The juristic method to use for calculating Asr timings.
	AsrMethod AsrJuristicMethod `sql:"type:ENUM('SHAFI_HANBALI_MALIKI','HANAFI')" gorm:"column:asr_method"`
	// The high latitude rule to use to calculate Fajr and Isha prayers.
	HighLatitudeRule HighLatitudeRule `sql:"type:ENUM('NO_HIGH_LATITUDE_RULE','MIDDLE_OF_THE_NIGHT','SEVENTH_OF_THE_NIGHT','TWILIGHT_ANGLE')" gorm:"column:high_latitude_rule"`
	// The prayer adjustments (aka offsets) to apply to the calculated prayer times.
	Adjustments PrayerAdjustments `gorm:"embedded"`
}
