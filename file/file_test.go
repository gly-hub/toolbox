package file

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestFile(t *testing.T) {
	pwd, err1 := GetPwd()
	assert.Nil(t, err1)
	assert.NotEqual(t, pwd, "")

	ok1, err2 := CheckFileIsExist(path.Join(pwd, "file.go"))
	assert.Nil(t, err2)
	assert.Equal(t, ok1, true)

	ok2, err3 := CheckFileIsExist("")
	assert.NotNil(t, err3)
	assert.NotEqual(t, ok2, true)

	err4 := CreateFile("test.txt")
	assert.Nil(t, err4)

	err5 := CreateFile("test.txt")
	assert.Nil(t, err5)

	err6 := CreateFile("")
	assert.NotNil(t, err6)

	err9 := WriteFile("test.txt", "test")
	assert.Nil(t, err9)

	err10 := WriteFile("", "test")
	assert.NotNil(t, err10)

	err7 := DeleteFile("")
	assert.Nil(t, err7)

	err8 := DeleteFile(path.Join(pwd, "test.txt"))
	assert.Nil(t, err8)
}
