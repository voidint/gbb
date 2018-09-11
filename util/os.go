package util

import (
	"fmt"
	"os"
)

// FileExist 判断指定路径的文件是否存在
func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return os.IsExist(err)
	}
	return !info.IsDir()
}

// Chdir 切换到指定目录
func Chdir(dir string, debug bool) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if wd == dir {
		return nil
	}

	if debug {
		fmt.Printf("==> cd %s\n", dir)
	}
	return os.Chdir(dir)
}
