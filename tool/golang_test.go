package tool

import (
	"testing"

	"os"

	"path/filepath"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHasMain(t *testing.T) {
	Convey("检测指定的go源代码文件中是否包含main函数", t, func() {
		Convey("不包含main函数的go源代码文件", func() {
			yes, err := hasMain("./golang.go")
			So(err, ShouldBeNil)
			So(yes, ShouldBeFalse)
		})
		Convey("包含main函数的go源代码文件", func() {
			yes, err := hasMain("../main.go")
			So(err, ShouldBeNil)
			So(yes, ShouldBeTrue)
		})
		Convey("非go源代码文件", func() {
			yes, err := hasMain("../gbb.json")
			So(err, ShouldNotBeNil)
			So(yes, ShouldBeFalse)
		})
	})
}

func TestWalkMainDir(t *testing.T) {
	Convey("遍历根目录及其子目录查找包含main函数的源代码文件路径", t, func() {
		dir, err := os.Getwd()
		So(err, ShouldBeNil)
		So(strings.HasSuffix(dir, filepath.Join("src", "github.com", "voidint", "gbb", "tool")), ShouldBeTrue)
		workspace := strings.TrimRight(dir, "tool")
		paths, err := walkMainDir(workspace)
		So(err, ShouldBeNil)
		So(paths, ShouldNotBeEmpty)
		So(len(paths), ShouldEqual, 1)
		So(paths[0], ShouldEqual, filepath.Join(workspace, "main.go"))
	})
}
