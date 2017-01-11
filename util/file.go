package util

import "os"

// FileExist 判断指定路径的文件是否存在
func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return os.IsExist(err)
	}
	return !info.IsDir()
}
