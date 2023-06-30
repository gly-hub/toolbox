package captcha

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRestoreAsset(t *testing.T) {
	err := RestoreAsset("", "test.ttf")
	assert.NotNil(t, err)

	err2 := RestoreAsset("", "default.ttf")
	assert.Nil(t, err2)
}
