package test_infra

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mnadev/limestone/proto"
)

type Resource struct {
	CreateTime *timestamppb.Timestamp
	UpdateTime *timestamppb.Timestamp
}

// AssertProtoEqual asserts that two protobufs are equal, ignoring the fields specified in `ignoreFields`.
// An empty proto needs to be passed in as the `typ` argument.
func AssertProtoEqual(t *testing.T, expected, actual, typ interface{}, ignore cmp.Option) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		assert.Fail(t, "expected is nil; cannot compare a nil value with a not nil value")
	}

	if actual == nil {
		assert.Fail(t, "actual is nil; cannot compare a nil value with a not nil value")
	}

	unexportOpts := cmpopts.IgnoreUnexported(typ)
	if !cmp.Equal(expected, actual, ignore, unexportOpts, protocmp.Transform()) {
		diff := cmp.Diff(expected, actual, ignore, unexportOpts, protocmp.Transform())
		assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s%s", expected, actual, diff))
	}
}

// AssertUserTimestampsCurrent asserts that the create and update timestamps of the user are current.
// By current, it means that the timestamps are within a range of 1 second from now).
func AssertUserTimestampsCurrent(t *testing.T, u *pb.User) {
	assert.LessOrEqual(t, time.Now().Unix()-u.CreateTime.GetSeconds(), int64(10))
	assert.LessOrEqual(t, time.Now().Unix()-u.UpdateTime.GetSeconds(), int64(10))
}

// AssertMasjidTimestampsCurrent asserts that the create and update timestamps of the masjid are current.
// By current, it means that the timestamps are within a range of 1 second from now).
func AssertMasjidTimestampsCurrent(t *testing.T, m *pb.Masjid) {
	assert.LessOrEqual(t, time.Now().Unix()-m.CreateTime.GetSeconds(), int64(10))
	assert.LessOrEqual(t, time.Now().Unix()-m.UpdateTime.GetSeconds(), int64(10))
}

// AssertEventTimestampsCurrent asserts that the create and update timestamps of the event are current.
// By current, it means that the timestamps are within a range of 1 second from now).
func AssertEventTimestampsCurrent(t *testing.T, e *pb.Event) {
	assert.LessOrEqual(t, time.Now().Unix()-e.CreateTime.GetSeconds(), int64(10))
	assert.LessOrEqual(t, time.Now().Unix()-e.UpdateTime.GetSeconds(), int64(10))
}

// AssertAdhanFileTimestampsCurrent asserts that the create and update timestamps of the event are current.
// By current, it means that the timestamps are within a range of 1 second from now).
func AssertAdhanFileTimestampsCurrent(t *testing.T, a *pb.AdhanFile) {
	assert.LessOrEqual(t, time.Now().Unix()-a.CreateTime.GetSeconds(), int64(10))
	assert.LessOrEqual(t, time.Now().Unix()-a.UpdateTime.GetSeconds(), int64(10))
}
