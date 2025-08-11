package entity

type Education int32

const (
	EducationUnspecified Education = iota
	EducationHighSchool
	EducationBachelor
	EducationMaster
	EducationDoctorate
)

func (e Education) String() string {
	switch e {
	case EducationHighSchool:
		return "HIGH_SCHOOL"
	case EducationBachelor:
		return "BACHELOR"
	case EducationMaster:
		return "MASTER"
	case EducationDoctorate:
		return "DOCTORATE"
	default:
		return "EDUCATION_UNSPECIFIED"
	}
} 