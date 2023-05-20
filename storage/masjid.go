package storage

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	mpb "github.com/mnadev/limestone/masjid_service/proto"
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
	PhoneCountryCode string
	Number           string
	Extension        string
}

type Masjid struct {
	ID           uuid.UUID                `gorm:"primaryKey;type:char(36)"`
	Name         string                   `gorm:"unique;type:varchar(320)"`
	IsVerified   bool                     `gorm:"default:false"`
	Address      Address                  `gorm:"embedded"`
	PhoneNumber  PhoneNumber              `gorm:"embedded"`
	PrayerConfig PrayerTimesConfiguration `gorm:"embedded"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
			PhoneCountryCode: m.GetPhoneNumber().GetCountryCode(),
			Number:           m.GetPhoneNumber().GetNumber(),
			Extension:        m.GetPhoneNumber().GetExtension(),
		},
		PrayerConfig: PrayerTimesConfiguration{
			CalculationMethod: FromMasjidToInternalCalculationMethodEnum(m.GetPrayerConfig().GetMethod()),
			FajrAngle:         m.GetPrayerConfig().GetFajrAngle(),
			IshaAngle:         m.GetPrayerConfig().GetIshaAngle(),
			IshaInterval:      m.GetPrayerConfig().GetIshaInterval(),
			AsrMethod:         FromMasjidToInternalAsrMethodEnum(m.GetPrayerConfig().GetAsrMethod()),
			HighLatitudeRule:  FromMasjidToInternalHighLatitudeEnum(m.GetPrayerConfig().GetHighLatitudeRule()),
			Adjustments: PrayerAdjustments{
				FajrAdjustment:    m.GetPrayerConfig().GetAdjustments().GetFajrAdjustment(),
				DhuhrAdjustment:   m.GetPrayerConfig().GetAdjustments().GetDhuhrAdjustment(),
				AsrAdjustment:     m.GetPrayerConfig().GetAdjustments().GetAsrAdjustment(),
				MaghribAdjustment: m.GetPrayerConfig().GetAdjustments().GetMaghribAdjustment(),
				IshaAdjustment:    m.GetPrayerConfig().GetAdjustments().GetIshaAdjustment(),
			},
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
			CountryCode: m.PhoneNumber.PhoneCountryCode,
			Number:      m.PhoneNumber.Number,
			Extension:   m.PhoneNumber.Extension,
		},
		PrayerConfig: &mpb.PrayerTimesConfiguration{
			Method:           FromInternalToMasjidCalculationMethodEnum(m.PrayerConfig.CalculationMethod),
			FajrAngle:        m.PrayerConfig.FajrAngle,
			IshaAngle:        m.PrayerConfig.IshaAngle,
			IshaInterval:     m.PrayerConfig.IshaInterval,
			AsrMethod:        FromInternalToMasjidAsrMethodEnum(m.PrayerConfig.AsrMethod),
			HighLatitudeRule: FromInternalToMasjidHighLatitudeEnum(m.PrayerConfig.HighLatitudeRule),
			Adjustments: &mpb.PrayerTimesConfiguration_PrayerAdjustments{
				FajrAdjustment:    m.PrayerConfig.Adjustments.FajrAdjustment,
				DhuhrAdjustment:   m.PrayerConfig.Adjustments.DhuhrAdjustment,
				AsrAdjustment:     m.PrayerConfig.Adjustments.AsrAdjustment,
				MaghribAdjustment: m.PrayerConfig.Adjustments.MaghribAdjustment,
				IshaAdjustment:    m.PrayerConfig.Adjustments.IshaAdjustment,
			},
		},
		CreateTime: timestamppb.New(m.CreatedAt),
		UpdateTime: timestamppb.New(m.UpdatedAt),
	}
}
