package storage

// BirthDate represents the birth date of a user.
type BirthDate struct {
	Year  int32 `gorm:"column:year"`
	Month Month `gorm:"column:month"`
	Day   int8  `gorm:"column:day"`
}

// Month defines an enum for the month of the year.
type Month int

const (
	// Default value.
	MonthUnspecified Month = iota
	MonthJanuary
	MonthFebruary
	MonthMarch
	MonthApril
	MonthMay
	MonthJune
	MonthJuly
	MonthAugust
	MonthSeptember
	MonthOctober
	MonthNovember
	MonthDecember
)
