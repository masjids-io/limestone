package test_infra

import (
	"testing"
	"time"

	userservicepb "github.com/mnadev/limestone/user_service/proto"
	"github.com/stretchr/testify/assert"
)

func AssertUserProtoEqual(t *testing.T, expected, actual *userservicepb.User) bool {
	actual_without_timestamp := userservicepb.User{
		UserId:      actual.UserId,
		Email:       actual.Email,
		Username:    actual.Username,
		FirstName:   actual.FirstName,
		LastName:    actual.LastName,
		PhoneNumber: actual.PhoneNumber,
		Gender:      actual.Gender,
	}

	assert.Equal(t, *expected, actual_without_timestamp)
	assert.LessOrEqual(t, time.Now().Unix()-actual.CreateTime.GetSeconds(), int64(1))
	assert.LessOrEqual(t, time.Now().Unix()-actual.UpdateTime.GetSeconds(), int64(1))

	return true
}
