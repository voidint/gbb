package tool

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bouk/monkey"
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

func TestWalkMainPkgs(t *testing.T) {
	Convey("遍历根目录及其子目录查找包含main函数的源代码文件路径", t, func() {
		dir, err := os.Getwd()
		So(err, ShouldBeNil)
		So(strings.HasSuffix(dir, filepath.Join("src", "github.com", "voidint", "gbb", "tool")), ShouldBeTrue)
		workspace := strings.TrimRight(dir, fmt.Sprintf("%ctool", os.PathSeparator))
		paths, err := walkMainPkgs(workspace)
		So(err, ShouldBeNil)
		So(paths, ShouldNotBeEmpty)
		So(len(paths), ShouldEqual, 1)
		So(paths[0], ShouldEqual, workspace)
	})
}

func TestWalkPkgs(t *testing.T) {
	Convey("查找指定目录及其子目录下所有满足golang包的目录路径", t, func() {
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		paths, err := walkPkgs(wd)
		So(err, ShouldBeNil)
		So(paths, ShouldNotBeEmpty)
		So(len(paths), ShouldEqual, 1)
		So(paths[0], ShouldEqual, wd)

		paths, err = walkPkgs(strings.TrimRight(wd, "tool"))
		So(err, ShouldBeNil)
		So(paths, ShouldNotBeEmpty)
		So(len(paths), ShouldEqual, 7)
		So(paths, ShouldContain, strings.TrimRight(wd, "tool"))
		So(paths, ShouldContain, filepath.Join(strings.TrimRight(wd, "tool"), "build"))
		So(paths, ShouldContain, filepath.Join(strings.TrimRight(wd, "tool"), "cmd"))
		So(paths, ShouldContain, filepath.Join(strings.TrimRight(wd, "tool"), "config"))
		So(paths, ShouldContain, filepath.Join(strings.TrimRight(wd, "tool"), "tool"))
		So(paths, ShouldContain, filepath.Join(strings.TrimRight(wd, "tool"), "util"))
		So(paths, ShouldContain, filepath.Join(strings.TrimRight(wd, "tool"), "variable"))

		Convey("检查指定路径是否是golang包路径报错", func() {
			var ErrIsGoPkg = errors.New("error for test")
			monkey.Patch(isGoPkg, func(path string) (yes bool, err error) {
				return false, ErrIsGoPkg
			})
			defer monkey.Unpatch(isGoPkg)

			paths, err := walkPkgs(wd)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrIsGoPkg)
			So(paths, ShouldBeEmpty)
		})
	})
}

func TestIsGoPkg(t *testing.T) {
	Convey("判断是否是golang包目录", t, func() {
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(wd, ShouldNotBeBlank)
		Convey("合法路径", func() {
			Convey("路径下包含的全部是go源文件", func() {
				yes, err := isGoPkg(wd)
				So(err, ShouldBeNil)
				So(yes, ShouldBeTrue)
			})
			Convey("路径下包含的全部是目录，不包含任何go源文件", func() {
				path := filepath.Join(wd, "test")
				So(os.MkdirAll(filepath.Join(wd, "test", "subtest0"), 0755), ShouldBeNil)
				So(os.MkdirAll(filepath.Join(wd, "test", "subtest1"), 0755), ShouldBeNil)
				defer os.RemoveAll(path)

				yes, err := isGoPkg(path)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})

			Convey("路径下既包含目录，还包含go源文件", func() {
				yes, err := isGoPkg(strings.TrimRight(wd, "tool"))
				So(err, ShouldBeNil)
				So(yes, ShouldBeTrue)
			})
		})
		Convey("非法路径", func() {
			Convey("路径为空", func() {
				yes, err := isGoPkg("")
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径非目录", func() {
				yes, err := isGoPkg(filepath.Join(wd, "golang_test.go"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径不存在", func() {
				yes, err := isGoPkg(filepath.Join(wd, "not_exist_dir"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
		})
	})
}

func TestIsMainPkg(t *testing.T) {
	Convey("检查指定路径是否是main包", t, func() {
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(wd, ShouldNotBeBlank)

		Convey("合法路径", func() {
			Convey("非main包且不包含子包", func() {
				yes, err := isMainPkg(wd)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("非main包且包含子包", func() {
				path := filepath.Join(
					strings.TrimRight(wd, "tool"),
					"vendor", "github.com", "smartystreets", "assertions",
				)
				yes, err := isMainPkg(path)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})

			Convey("main包", func() {
				yes, err := isMainPkg(strings.TrimRight(wd, "tool"))
				So(err, ShouldBeNil)
				So(yes, ShouldBeTrue)
			})
		})

		Convey("非法路径", func() {
			Convey("路径为空", func() {
				yes, err := isMainPkg("")
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径非目录", func() {
				yes, err := isMainPkg(filepath.Join(wd, "golang_test.go"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径不存在", func() {
				yes, err := isMainPkg(filepath.Join(wd, "not_exist_dir"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径下不包含go源文件(非包路径)", func() {
				path := filepath.Join(
					strings.TrimRight(wd, "tool"),
					"vendor", "github.com",
				)
				yes, err := isMainPkg(path)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
		})
	})
}
