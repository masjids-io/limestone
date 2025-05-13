package entity

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
)

func (g Gender) String() string {
	if g == Male {
		return "MALE"
	}
	return "FEMALE"
}
