package entity

type Sect int32

const (
	SectUnspecified Sect = iota
	SectSunni
	SectSunniHanafi
	SectSunniMaleki
	SectSunniShafii
	SectSunniMaliki
	SectShia
	SectOther
)

func (s Sect) String() string {
	switch s {
	case SectSunni:
		return "SUNNI"
	case SectSunniHanafi:
		return "SUNNI_HANAFI"
	case SectSunniMaleki:
		return "SUNNI_MALEKI"
	case SectSunniShafii:
		return "SUNNI_SHAFII"
	case SectSunniMaliki:
		return "SUNNI_MALIKI"
	case SectShia:
		return "SHIA"
	case SectOther:
		return "OTHER"
	default:
		return "SECT_UNSPECIFIED"
	}
} 