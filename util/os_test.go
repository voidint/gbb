package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func getFilename() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
}

func TestFileExist(t *testing.T) {
	Convey("检测文件是否存在", t, func() {
		Convey("预期文件存在，实际文件存在", func() {
			filename := getFilename()
			err := ioutil.WriteFile(filename, []byte("hello world"), 0777)
			So(err, ShouldBeNil)
			defer os.Remove(filename)
			So(FileExist(filename), ShouldBeTrue)
		})

		Convey("预期文件存在，实际文件不存在", func() {
			Convey("文件和文件夹都不存在", func() {
				filename := getFilename()
				So(FileExist(filename), ShouldBeFalse)
			})
			Convey("文件不存在，文件夹存在", func() {
				filename := getFilename()
				os.MkdirAll(filename, 0666)
				So(FileExist(filename), ShouldBeFalse)
			})
		})
	})
}

func TestChdir(t *testing.T) {
	Convey("切换工作目录", t, func() {
		Convey("目标目录是当前目录", func() {
			wd, err := os.Getwd()
			So(err, ShouldBeNil)
			So(Chdir(wd, true), ShouldBeNil)
		})

		Convey("目标目录非当前目录", func() {
			wd, err := os.Getwd()
			So(err, ShouldBeNil)

			defer Chdir(wd, true) // init work directory

			if idx := strings.LastIndex(wd, fmt.Sprintf("%c", os.PathSeparator)); idx > 0 {
				So(Chdir(wd[:idx], true), ShouldBeNil)
			}
		})

		Convey("目录切换发生错误", func() {
			var ErrChdir = errors.New("chdir error")
			monkey.Patch(os.Getwd, func() (dir string, err error) {
				return "", ErrChdir
			})
			defer monkey.Unpatch(os.Getwd)
			So(Chdir("../", true), ShouldEqual, ErrChdir)
		})
	})
}
