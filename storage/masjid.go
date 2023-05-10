package storage

import (
	"time"

	"github.com/google/uuid"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Address struct {
	AddressLine1 string
	AddressLine2 string
	ZoneCode     string
	PostalCode   string
	City         string
	CountryCode  string
}

type PhoneNumber struct {
	CountryCode string
	Number      string
	Extension   string
}

type Masjid struct {
	ID          uuid.UUID   `gorm:"primaryKey;type:char(36)"`
	Name        string      `gorm:"unique;type:varchar(320)"`
	IsVerified  bool        `gorm:"default:false"`
	Address     Address     `gorm:"embedded"`
	PhoneNumber PhoneNumber `gorm:"embedded"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewMasjid creates a new Masjid struct given the Masjid proto.
func NewMasjid(m *mpb.Masjid) (*Masjid, error) {
	return &Masjid{
		Name:       m.GetName(),
		IsVerified: m.GetIsVerified(),
		Address: Address{
			AddressLine1: m.GetAddress().GetAddressLine_1(),
			AddressLine2: m.GetAddress().GetAddressLine_2(),
			ZoneCode:     m.GetAddress().GetZoneCode(),
			PostalCode:   m.GetAddress().GetPostalCode(),
			City:         m.GetAddress().GetCity(),
			CountryCode:  m.GetAddress().GetCountryCode(),
		},
		PhoneNumber: PhoneNumber{
			CountryCode: m.GetPhoneNumber().GetCountryCode(),
			Number:      m.GetPhoneNumber().GetNumber(),
			Extension:   m.GetPhoneNumber().GetExtension(),
		},
	}, nil
}

func (m *Masjid) ToProto() *mpb.Masjid {
	return &mpb.Masjid{
		Id:         m.ID.String(),
		Name:       m.Name,
		IsVerified: m.IsVerified,
		Address: &mpb.Masjid_Address{
			AddressLine_1: m.Address.AddressLine1,
			AddressLine_2: m.Address.AddressLine2,
			ZoneCode:      m.Address.ZoneCode,
			PostalCode:    m.Address.PostalCode,
			City:          m.Address.City,
			CountryCode:   m.Address.CountryCode,
		},
		PhoneNumber: &mpb.Masjid_PhoneNumber{
			CountryCode: m.PhoneNumber.CountryCode,
			Number:      m.PhoneNumber.Number,
			Extension:   m.PhoneNumber.Extension,
		},
		CreateTime: timestamppb.New(m.CreatedAt),
		UpdateTime: timestamppb.New(m.UpdatedAt),
	}
}
