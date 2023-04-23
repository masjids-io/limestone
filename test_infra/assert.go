package test_infra

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	userservicepb "github.com/mnadev/limestone/user_service/proto"
	"github.com/stretchr/testify/assert"
)

// AssertProtoEqual asserts that two protobufs are equal, ignoring the fields specified in `ignoreFields`.
// An empty proto needs to be passed in as the `typ` argument.
func AssertProtoEqual(t *testing.T, expected, actual, typ interface{}, ignoreFields ...string) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		assert.Fail(t, "expected is nil; cannot compare a nil value with a not nil value")
	}

	if actual == nil {
		assert.Fail(t, "actual is nil; cannot compare a nil value with a not nil value")
	}

	fieldOpts := cmpopts.IgnoreFields(typ, ignoreFields...)
	unexportOpts := cmpopts.IgnoreUnexported(typ)
	if !cmp.Equal(expected, actual, fieldOpts, unexportOpts) {
		diff := cmp.Diff(expected, actual, fieldOpts, unexportOpts)
		assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s%s", expected, actual, diff))
	}
}

// AssertTimestampsCurrent asserts that the create and update timestamps of the user are current.
// By current, it means that the timestamps are within a range of 1 second from now).
func AssertTimestampsCurrent(t *testing.T, u *userservicepb.User) {
	assert.LessOrEqual(t, time.Now().Unix()-u.CreateTime.GetSeconds(), int64(1))
	assert.LessOrEqual(t, time.Now().Unix()-u.UpdateTime.GetSeconds(), int64(1))
}
