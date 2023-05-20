package storage

import (
	"time"

	"github.com/google/uuid"
	epb "github.com/mnadev/limestone/event_service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GenderRestriction int64

const (
	NO_RESTRICTION GenderRestriction = iota
	MALE_ONLY
	FEMALE_ONLY
)

type EventType int64

const (
	WORSHIP EventType = iota
	EDUCATIONAL
	COMMUNITY
	ATHLETIC
	FUNDRAISING
	YOUTH
	CHILDREN_SPECIFIC
	MATRIMONIAL
	FUNERAL
)

type Event struct {
	ID                uuid.UUID `gorm:"primaryKey;type:char(36)"`
	UserId            string
	MasjidId          string
	Name              string
	Description       string
	StartTime         time.Time
	EndTime           time.Time
	GenderRestriction GenderRestriction `sql:"type:ENUM('NO_RESTRICTION', 
														'MALE_ONLY',
														'FEMALE_ONLY')" 
														gorm:"column:gender_restriction"`
	// The different event types associated with this event.
	// EventType       EventType
	IsPaid          bool
	RequiresRsvp    bool
	MaxParticipants int32
	LivestreamLink  string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewEvent creates a new Event struct given the Event proto.
func NewEvent(ep *epb.Event) (*Event, error) {
	e := Event{
		Name:              ep.GetName(),
		Description:       ep.GetDescription(),
		StartTime:         ep.GetStartTime().AsTime(),
		EndTime:           ep.GetEndTime().AsTime(),
		GenderRestriction: FromProtoToInternalGenderRestriction(ep.GetGenderRestriction()),
		IsPaid:            ep.GetIsPaid(),
		RequiresRsvp:      ep.GetRequiresRsvp(),
		MaxParticipants:   ep.GetMaxParticipants(),
		LivestreamLink:    ep.GetLivestreamLink(),
	}
	if ep.GetUserId() != "" {
		e.UserId = ep.GetUserId()
	} else {
		e.MasjidId = ep.GetMasjidId()
	}

	return &e, nil
}

func (e *Event) ToProto() *epb.Event {
	ep := epb.Event{
		Name:              e.Name,
		Description:       e.Description,
		StartTime:         timestamppb.New(e.StartTime),
		EndTime:           timestamppb.New(e.EndTime),
		GenderRestriction: FromInternalToProtoGenderRestriction(e.GenderRestriction),
		IsPaid:            e.IsPaid,
		RequiresRsvp:      e.RequiresRsvp,
		MaxParticipants:   e.MaxParticipants,
		LivestreamLink:    e.LivestreamLink,
		CreateTime:        timestamppb.New(e.CreatedAt),
		UpdateTime:        timestamppb.New(e.UpdatedAt),
	}
	if e.UserId != "" {
		ep.Owner = &epb.Event_UserId{UserId: e.UserId}
	} else {
		ep.Owner = &epb.Event_MasjidId{MasjidId: e.MasjidId}
	}

	return &ep
}

func FromProtoToInternalGenderRestriction(g epb.Event_GenderRestriction) GenderRestriction {
	switch g {
	case epb.Event_MALE_ONLY:
		return MALE_ONLY
	case epb.Event_FEMALE_ONLY:
		return FEMALE_ONLY
	}
	return NO_RESTRICTION
}

func FromInternalToProtoGenderRestriction(g GenderRestriction) epb.Event_GenderRestriction {
	switch g {
	case MALE_ONLY:
		return epb.Event_MALE_ONLY
	case FEMALE_ONLY:
		return epb.Event_FEMALE_ONLY
	}
	return epb.Event_NO_RESTRICTION
}
