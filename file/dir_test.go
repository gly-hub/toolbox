package file

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestDir(t *testing.T) {
	pwd, _ := GetPwd()
	ok1, err1 := CheckDirExist(pwd)
	assert.Nil(t, err1)
	assert.Equal(t, ok1, true)

	ok2, err2 := CheckDirExist(path.Join(pwd, "test"))
	assert.NotNil(t, err2)
	assert.Equal(t, ok2, false)

	err3 := CreateDir(path.Join(pwd, "test"))
	assert.Nil(t, err3)

	err4 := DeleteDir(path.Join(pwd, "test"))
	assert.Nil(t, err4)
}
