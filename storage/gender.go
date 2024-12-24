package storage

type gender string

const (
	MALE   gender = "MALE"
	FEMALE gender = "FEMALE"
)

func (g gender) String() string {
	if g == MALE {
		return "MALE"
	}
	return "FEMALE"
}
