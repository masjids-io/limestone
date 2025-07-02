package entity

type Gender string

const (
	Male              Gender = "MALE"
	Female            Gender = "FEMALE"
	GenderUnspecified Gender = "UNSPECIFIED"
)

func (g Gender) String() string {
	if g == Male {
		return "MALE"
	} else if g == Female {
		return "FEMALE"
	}
	return "UNSPECIFIED"
}
