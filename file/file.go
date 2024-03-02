package file

import (
	"os"
)

// GetPwd 获取当前目录
func GetPwd() (string, error) {
	return os.Getwd()
}

// CheckFileIsExist 校验文件是否存在
func CheckFileIsExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}

	if !fileInfo.IsDir() {
		return true, nil
	}

	return false, nil
}

// CreateFile 创建文件
func CreateFile(path string) error {
	if ok, _ := CheckFileIsExist(path); ok {
		return nil
	}

	_, err := os.Create(path)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile 删除文件
func DeleteFile(path string) error {
	if ok, _ := CheckFileIsExist(path); !ok {
		return nil
	}

	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// WriteFile 写入文件
func WriteFile(path string, content string) error {

	err := os.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
