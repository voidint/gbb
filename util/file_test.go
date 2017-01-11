package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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
