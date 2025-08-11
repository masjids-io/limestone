package entity

type Location struct {
	Country   string `gorm:"type:varchar(255)"`
	City      string `gorm:"type:varchar(255)"`
	State     string `gorm:"type:varchar(255)"`
	ZipCode   string `gorm:"type:varchar(20)"`
	Latitude  int32  `gorm:"type:int"`
	Longitude int32  `gorm:"type:int"`
} 