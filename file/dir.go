package file

import "os"

// CheckDirExist 校验文件夹是否存在
func CheckDirExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, err
	}

	if fileInfo.IsDir() {
		return true, nil
	}

	return false, nil
}

// CreateDir 创建文件夹
func CreateDir(path string) error {
	if ok, _ := CheckDirExist(path); ok {
		return nil
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDir 删除文件夹
func DeleteDir(path string) error {
	if ok, _ := CheckDirExist(path); !ok {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
