package entity

type Hobbies int32

const (
	HobbiesUnspecified Hobbies = iota
	HobbiesReading
	HobbiesWatchingMovies
	HobbiesWatchingTV
	HobbiesPainting
	HobbiesWriting
	HobbiesVolunteering
	HobbiesVolleyball
	HobbiesTraveling
	HobbiesGaming
	HobbiesMartialArts
	HobbiesSoccer
	HobbiesBasketball
	HobbiesFootball
	HobbiesTennis
	HobbiesBadminton
	HobbiesCricket
)

func (h Hobbies) String() string {
	switch h {
	case HobbiesReading:
		return "READING"
	case HobbiesWatchingMovies:
		return "WATCHING_MOVIES"
	case HobbiesWatchingTV:
		return "WATCHING_TV"
	case HobbiesPainting:
		return "PAINTING"
	case HobbiesWriting:
		return "WRITING"
	case HobbiesVolunteering:
		return "VOLUNTEERING"
	case HobbiesVolleyball:
		return "VOLLEYBALL"
	case HobbiesTraveling:
		return "TRAVELING"
	case HobbiesGaming:
		return "GAMING"
	case HobbiesMartialArts:
		return "MARTIAL_ARTS"
	case HobbiesSoccer:
		return "SOCCER"
	case HobbiesBasketball:
		return "BASKETBALL"
	case HobbiesFootball:
		return "FOOTBALL"
	case HobbiesTennis:
		return "TENNIS"
	case HobbiesBadminton:
		return "BADMINTON"
	case HobbiesCricket:
		return "CRICKET"
	default:
		return "HOBBIES_UNSPECIFIED"
	}
} 