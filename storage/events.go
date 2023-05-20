package storage

import (
	"strings"
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
	EventTypes      string
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

	types := []string{}

	for _, t := range ep.Types {
		types = append(types, FromProtoToInternalEventType(t))
	}

	e.EventTypes = strings.Join(types, ",")

	return &e, nil
}

func (e *Event) ToProto() *epb.Event {
	ep := epb.Event{
		Id:                e.ID.String(),
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

	types := e.EventTypes
	typespb := []epb.Event_EventType{}
	for _, t := range strings.Split(types, ",") {
		typespb = append(typespb, FromInternalToProtoEvent(t))
	}
	ep.Types = typespb

	return &ep
}

func FromProtoToInternalEventType(et epb.Event_EventType) string {
	switch et {
	case epb.Event_EDUCATIONAL:
		return "EDUCATIONAL"
	case epb.Event_COMMUNITY:
		return "COMMUNITY"
	case epb.Event_ATHLETIC:
		return "ATHLETIC"
	case epb.Event_FUNDRAISING:
		return "FUNDRAISING"
	case epb.Event_YOUTH:
		return "YOUTH"
	case epb.Event_CHILDREN_SPECIFIC:
		return "CHILDREN_SPECIFIC"
	case epb.Event_MATRIMONIAL:
		return "MATRIMONIAL"
	case epb.Event_FUNERAL:
		return "FUNERAL"
	case epb.Event_WORSHIP:
		return "WORSHIP"
	}
	return ""
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

func FromInternalToProtoEvent(s string) epb.Event_EventType {
	switch s {
	case "EDUCATIONAL":
		return epb.Event_EDUCATIONAL
	case "COMMUNITY":
		return epb.Event_COMMUNITY
	case "ATHLETIC":
		return epb.Event_ATHLETIC
	case "FUNDRAISING":
		return epb.Event_FUNDRAISING
	case "YOUTH":
		return epb.Event_YOUTH
	case "CHILDREN_SPECIFIC":
		return epb.Event_CHILDREN_SPECIFIC
	case "MATRIMONIAL":
		return epb.Event_MATRIMONIAL
	case "FUNERAL":
		return epb.Event_FUNERAL
	}
	return epb.Event_WORSHIP
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
