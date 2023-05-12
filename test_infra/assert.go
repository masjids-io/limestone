package test_infra

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	upb "github.com/mnadev/limestone/user_service/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func AssertUserTimestampsCurrent(t *testing.T, u *upb.User) {
	assert.LessOrEqual(t, time.Now().Unix()-u.CreateTime.GetSeconds(), int64(1))
	assert.LessOrEqual(t, time.Now().Unix()-u.UpdateTime.GetSeconds(), int64(1))
}

// AssertMasjidTimestampsCurrent asserts that the create and update timestamps of the masjid are current.
// By current, it means that the timestamps are within a range of 1 second from now).
func AssertMasjidTimestampsCurrent(t *testing.T, m *mpb.Masjid) {
	assert.LessOrEqual(t, time.Now().Unix()-m.CreateTime.GetSeconds(), int64(1))
	assert.LessOrEqual(t, time.Now().Unix()-m.UpdateTime.GetSeconds(), int64(1))
}
