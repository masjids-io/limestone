package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHash_Success(t *testing.T) {
	p := "password"

	h, err := HashPassword(p)
	assert.Nil(t, err)
	assert.Len(t, h, 60)
}

func TestCheckPassword_Mismatch(t *testing.T) {
	p1 := "password1"

	h, err := HashPassword(p1)
	assert.Nil(t, err)

	p2 := "password2"
	err = CheckPassword(p2, h)
	assert.Error(t, err)
}

func TestCheckPassword_Succeeds(t *testing.T) {
	p := "password"

	h, err := HashPassword(p)
	assert.Nil(t, err)

	err = CheckPassword(p, h)
	assert.Nil(t, err)
}
