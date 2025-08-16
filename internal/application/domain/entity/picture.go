package entity

type Picture struct {
	Image    []byte `gorm:"type:bytea"`
	MimeType string `gorm:"type:varchar(100)"`
} 