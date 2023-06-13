package storage

import (
	pb "github.com/mnadev/limestone/proto"
)

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
	CalculationMethod CalculationMethod `sql:"type:ENUM('OTHER', 
														'MUSLIM_WORLD_LEAGUE', 
														'EGYPTIAN', 
														'KARACHI', 
														'UMM_AL_QURA', 
														'DUBAI', 
														'MOON_SIGHTING_COMMITTEE', 
														'NORTH_AMERICA', 
														'KUWAIT', 
														'QATAR', 
														'SINGAPORE', 
														'UOIF')" 
														gorm:"column:calculation_method"`
	// The Fajr angle to use. This does not need to be set if the calculation method is set.
	// If both fields are set, then the calculation method field takes precedence.
	FajrAngle float64 `gorm:"default:0"`
	// The Isha angle to use. This does not need to be set if the calculation method is set.
	// If both fields are set, then the calculation method field takes precedence.
	IshaAngle float64 `gorm:"default:0"`
	// Minutes after Maghrib (if set, the time for Isha will be Maghrib plus the Isha interval).
	IshaInterval int32 `gorm:"default:0"`
	// The juristic method to use for calculating Asr timings.
	AsrMethod AsrJuristicMethod `sql:"type:ENUM('SHAFI_HANBALI_MALIKI', 
												'HANAFI')" 
												gorm:"column:asr_method"`
	// The high latitude rule to use to calculate Fajr and Isha prayers.
	HighLatitudeRule HighLatitudeRule `sql:"type:ENUM('NO_HIGH_LATITUDE_RULE', 
													  'MIDDLE_OF_THE_NIGHT', 
													  'SEVENTH_OF_THE_NIGHT', 
													  'TWILIGHT_ANGLE')" 
													   gorm:"column:high_latitude_rule"`
	// The prayer adjusments (aka offsets) to apply to the calculated prayer times.
	Adjustments PrayerAdjustments `gorm:"embedded"`
}

func FromMasjidToInternalCalculationMethodEnum(c pb.PrayerTimesConfiguration_CalculationMethod) CalculationMethod {
	switch c {
	case pb.PrayerTimesConfiguration_OTHER:
		return OTHER
	case pb.PrayerTimesConfiguration_MUSLIM_WORLD_LEAGUE:
		return MUSLIM_WORLD_LEAGUE
	case pb.PrayerTimesConfiguration_EGYPTIAN:
		return EGYPTIAN
	case pb.PrayerTimesConfiguration_KARACHI:
		return KARACHI
	case pb.PrayerTimesConfiguration_UMM_AL_QURA:
		return UMM_AL_QURA
	case pb.PrayerTimesConfiguration_DUBAI:
		return DUBAI
	case pb.PrayerTimesConfiguration_MOON_SIGHTING_COMMITTEE:
		return MOON_SIGHTING_COMMITTEE
	case pb.PrayerTimesConfiguration_NORTH_AMERICA:
		return NORTH_AMERICA
	case pb.PrayerTimesConfiguration_KUWAIT:
		return KUWAIT
	case pb.PrayerTimesConfiguration_QATAR:
		return QATAR
	case pb.PrayerTimesConfiguration_SINGAPORE:
		return SINGAPORE
	case pb.PrayerTimesConfiguration_UOIF:
		return UOIF
	}
	return OTHER
}

func FromMasjidToInternalAsrMethodEnum(a pb.PrayerTimesConfiguration_AsrJuristicMethod) AsrJuristicMethod {
	switch a {
	case pb.PrayerTimesConfiguration_SHAFI_HANBALI_MALIKI:
		return SHAFI_HANBALI_MALIKI
	case pb.PrayerTimesConfiguration_HANAFI:
		return HANAFI
	}
	return SHAFI_HANBALI_MALIKI
}

func FromMasjidToInternalHighLatitudeEnum(h pb.PrayerTimesConfiguration_HighLatitudeRule) HighLatitudeRule {
	switch h {
	case pb.PrayerTimesConfiguration_NO_HIGH_LATITUDE_RULE:
		return NO_HIGH_LATITUDE_RULE
	case pb.PrayerTimesConfiguration_MIDDLE_OF_THE_NIGHT:
		return MIDDLE_OF_THE_NIGHT
	case pb.PrayerTimesConfiguration_SEVENTH_OF_THE_NIGHT:
		return SEVENTH_OF_THE_NIGHT
	case pb.PrayerTimesConfiguration_TWILIGHT_ANGLE:
		return TWILIGHT_ANGLE
	}
	return NO_HIGH_LATITUDE_RULE
}

func FromInternalToMasjidCalculationMethodEnum(c CalculationMethod) pb.PrayerTimesConfiguration_CalculationMethod {
	switch c {
	case OTHER:
		return pb.PrayerTimesConfiguration_OTHER
	case MUSLIM_WORLD_LEAGUE:
		return pb.PrayerTimesConfiguration_MUSLIM_WORLD_LEAGUE
	case EGYPTIAN:
		return pb.PrayerTimesConfiguration_EGYPTIAN
	case KARACHI:
		return pb.PrayerTimesConfiguration_KARACHI
	case UMM_AL_QURA:
		return pb.PrayerTimesConfiguration_UMM_AL_QURA
	case DUBAI:
		return pb.PrayerTimesConfiguration_DUBAI
	case MOON_SIGHTING_COMMITTEE:
		return pb.PrayerTimesConfiguration_MOON_SIGHTING_COMMITTEE
	case NORTH_AMERICA:
		return pb.PrayerTimesConfiguration_NORTH_AMERICA
	case KUWAIT:
		return pb.PrayerTimesConfiguration_KUWAIT
	case QATAR:
		return pb.PrayerTimesConfiguration_QATAR
	case SINGAPORE:
		return pb.PrayerTimesConfiguration_SINGAPORE
	case UOIF:
		return pb.PrayerTimesConfiguration_UOIF
	}
	return pb.PrayerTimesConfiguration_OTHER
}

func FromInternalToMasjidAsrMethodEnum(a AsrJuristicMethod) pb.PrayerTimesConfiguration_AsrJuristicMethod {
	switch a {
	case SHAFI_HANBALI_MALIKI:
		return pb.PrayerTimesConfiguration_SHAFI_HANBALI_MALIKI
	case HANAFI:
		return pb.PrayerTimesConfiguration_HANAFI
	}
	return pb.PrayerTimesConfiguration_SHAFI_HANBALI_MALIKI
}

func FromInternalToMasjidHighLatitudeEnum(h HighLatitudeRule) pb.PrayerTimesConfiguration_HighLatitudeRule {
	switch h {
	case NO_HIGH_LATITUDE_RULE:
		return pb.PrayerTimesConfiguration_NO_HIGH_LATITUDE_RULE
	case MIDDLE_OF_THE_NIGHT:
		return pb.PrayerTimesConfiguration_MIDDLE_OF_THE_NIGHT
	case SEVENTH_OF_THE_NIGHT:
		return pb.PrayerTimesConfiguration_SEVENTH_OF_THE_NIGHT
	case TWILIGHT_ANGLE:
		return pb.PrayerTimesConfiguration_TWILIGHT_ANGLE
	}
	return pb.PrayerTimesConfiguration_NO_HIGH_LATITUDE_RULE
}
