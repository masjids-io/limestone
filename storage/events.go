package storage

import (
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mnadev/limestone/proto"
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
func NewEvent(ep *pb.Event) (*Event, error) {
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
	} else if ep.GetMasjidId() != "" {
		e.MasjidId = ep.GetMasjidId()
	} else {
		return nil, status.Error(codes.InvalidArgument, "User ID or Masjid must be specified")
	}

	types := []string{}

	for _, t := range ep.Types {
		types = append(types, FromProtoToInternalEventType(t))
	}

	e.EventTypes = strings.Join(types, ",")

	return &e, status.Error(codes.OK, codes.OK.String())
}

func (e *Event) ToProto() *pb.Event {
	ep := pb.Event{
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
		ep.Owner = &pb.Event_UserId{UserId: e.UserId}
	} else if e.MasjidId != "" {
		ep.Owner = &pb.Event_MasjidId{MasjidId: e.MasjidId}
	} else {
		// TODO: add logs
	}

	types := e.EventTypes
	typespb := []pb.Event_EventType{}
	for _, t := range strings.Split(types, ",") {
		typespb = append(typespb, FromInternalToProtoEvent(t))
	}
	ep.Types = typespb

	return &ep
}

func FromProtoToInternalEventType(et pb.Event_EventType) string {
	switch et {
	case pb.Event_EDUCATIONAL:
		return "EDUCATIONAL"
	case pb.Event_COMMUNITY:
		return "COMMUNITY"
	case pb.Event_ATHLETIC:
		return "ATHLETIC"
	case pb.Event_FUNDRAISING:
		return "FUNDRAISING"
	case pb.Event_YOUTH:
		return "YOUTH"
	case pb.Event_CHILDREN_SPECIFIC:
		return "CHILDREN_SPECIFIC"
	case pb.Event_MATRIMONIAL:
		return "MATRIMONIAL"
	case pb.Event_FUNERAL:
		return "FUNERAL"
	case pb.Event_WORSHIP:
		return "WORSHIP"
	}
	return ""
}

func FromProtoToInternalGenderRestriction(g pb.Event_GenderRestriction) GenderRestriction {
	switch g {
	case pb.Event_MALE_ONLY:
		return MALE_ONLY
	case pb.Event_FEMALE_ONLY:
		return FEMALE_ONLY
	}
	return NO_RESTRICTION
}

func FromInternalToProtoEvent(s string) pb.Event_EventType {
	switch s {
	case "EDUCATIONAL":
		return pb.Event_EDUCATIONAL
	case "COMMUNITY":
		return pb.Event_COMMUNITY
	case "ATHLETIC":
		return pb.Event_ATHLETIC
	case "FUNDRAISING":
		return pb.Event_FUNDRAISING
	case "YOUTH":
		return pb.Event_YOUTH
	case "CHILDREN_SPECIFIC":
		return pb.Event_CHILDREN_SPECIFIC
	case "MATRIMONIAL":
		return pb.Event_MATRIMONIAL
	case "FUNERAL":
		return pb.Event_FUNERAL
	}
	return pb.Event_WORSHIP
}

func FromInternalToProtoGenderRestriction(g GenderRestriction) pb.Event_GenderRestriction {
	switch g {
	case MALE_ONLY:
		return pb.Event_MALE_ONLY
	case FEMALE_ONLY:
		return pb.Event_FEMALE_ONLY
	}
	return pb.Event_NO_RESTRICTION
}
