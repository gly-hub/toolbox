package ip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLocation(t *testing.T) {
	local1 := GetLocation("127.0.0.1")
	assert.Equal(t, local1, "内部IP")

	local2 := GetLocation("localhost")
	assert.Equal(t, local2, "内部IP")

	local3 := GetLocation("110.242.68.66")
	assert.NotEqual(t, local3, "内部IP")
}

func TestGetLocalHost(t *testing.T) {
	ip := GetLocalHost()
	assert.NotEqual(t, ip, "")
}
