package entity

type Role string

const (
	ROLE_UNSPECIFIED Role = "UNSPECIFIED"
	MASJID_MEMBER    Role = "MASJID_MEMBER"
	MASJID_VOLUNTEER Role = "MASJID_VOLUNTEER"
	MASJID_ADMIN     Role = "MASJID_ADMIN"
	MASJID_IMAM      Role = "MASJID_IMAM"
)

func (r Role) String() string {
	switch r {
	case ROLE_UNSPECIFIED:
		return "UNSPECIFIED"
	case MASJID_MEMBER:
		return "MASJID_MEMBER"
	case MASJID_VOLUNTEER:
		return "MASJID_VOLUNTEER"
	case MASJID_ADMIN:
		return "MASJID_ADMIN"
	case MASJID_IMAM:
		return "MASJID_IMAM"
	default:
		return "UNSPECIFIED"
	}
}
