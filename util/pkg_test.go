package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWalkPkgsFunc(t *testing.T) {
	Convey("返回指定目录下满足过滤条件的go package路径列表", t, func() {
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(strings.HasSuffix(wd, "util"), ShouldBeTrue)

		paths, _, err := WalkPkgsFunc(wd, IsGoPkg)
		So(err, ShouldBeNil)
		So(paths, ShouldNotBeEmpty)
		So(len(paths), ShouldEqual, 1)
		So(paths[0], ShouldEqual, wd)

		paths, _, err = WalkPkgsFunc(strings.TrimSuffix(wd, "util"), IsGoPkg)
		So(err, ShouldBeNil)
		So(paths, ShouldNotBeEmpty)
		So(len(paths), ShouldEqual, 7)
		So(paths, ShouldContain, strings.TrimSuffix(wd, "util"))
		So(paths, ShouldContain, filepath.Join(strings.TrimSuffix(wd, "util"), "build"))
		So(paths, ShouldContain, filepath.Join(strings.TrimSuffix(wd, "util"), "cmd"))
		So(paths, ShouldContain, filepath.Join(strings.TrimSuffix(wd, "util"), "config"))
		So(paths, ShouldContain, filepath.Join(strings.TrimSuffix(wd, "util"), "tool"))
		So(paths, ShouldContain, filepath.Join(strings.TrimSuffix(wd, "util"), "util"))
		So(paths, ShouldContain, filepath.Join(strings.TrimSuffix(wd, "util"), "variable"))

		Convey("检查指定路径是否是golang包路径报错", func() {
			var ErrIsGoPkg = errors.New("error for test")
			paths, _, err := WalkPkgsFunc(wd, func(_ string) (bool, error) {
				return false, ErrIsGoPkg
			})
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
		So(strings.HasSuffix(wd, "util"), ShouldBeTrue)

		Convey("合法路径", func() {
			Convey("路径下包含的全部是go源文件", func() {
				yes, err := IsGoPkg(wd)
				So(err, ShouldBeNil)
				So(yes, ShouldBeTrue)
			})
			Convey("路径下包含的全部是目录，不包含任何go源文件", func() {
				path := filepath.Join(wd, "test")
				So(os.MkdirAll(filepath.Join(wd, "test", "subtest0"), 0755), ShouldBeNil)
				So(os.MkdirAll(filepath.Join(wd, "test", "subtest1"), 0755), ShouldBeNil)
				defer os.RemoveAll(path)

				yes, err := IsGoPkg(path)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})

			Convey("路径下既包含目录，还包含go源文件", func() {
				yes, err := IsGoPkg(strings.TrimSuffix(wd, "util"))
				So(err, ShouldBeNil)
				So(yes, ShouldBeTrue)
			})
		})
		Convey("非法路径", func() {
			Convey("路径为空", func() {
				yes, err := IsGoPkg("")
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径非目录", func() {
				yes, err := IsGoPkg(filepath.Join(wd, "pkg_test.go"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径不存在", func() {
				yes, err := IsGoPkg(filepath.Join(wd, "not_exist_dir"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
		})
	})
}

func TestIsMainPkg(t *testing.T) {
	Convey("判断指定路径是否是golang main package", t, func() {
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(wd, ShouldNotBeBlank)
		So(strings.HasSuffix(wd, "util"), ShouldBeTrue)

		Convey("合法路径", func() {
			Convey("路径下包含的全部是go源文件", func() {
				yes, err := IsMainPkg(wd)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径下包含的全部是目录，不包含任何go源文件", func() {
				path := filepath.Join(wd, "test")
				So(os.MkdirAll(filepath.Join(wd, "test", "subtest0"), 0755), ShouldBeNil)
				So(os.MkdirAll(filepath.Join(wd, "test", "subtest1"), 0755), ShouldBeNil)
				defer os.RemoveAll(path)

				yes, err := IsMainPkg(path)
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})

			Convey("路径下既包含目录，还包含go源文件", func() {
				yes, err := IsMainPkg(strings.TrimSuffix(wd, "util"))
				So(err, ShouldBeNil)
				So(yes, ShouldBeTrue)
			})
		})
		Convey("非法路径", func() {
			Convey("路径为空", func() {
				yes, err := IsMainPkg("")
				So(err, ShouldBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径非目录", func() {
				yes, err := IsMainPkg(filepath.Join(wd, "pkg_test.go"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
			Convey("路径不存在", func() {
				yes, err := IsMainPkg(filepath.Join(wd, "not_exist_dir"))
				So(err, ShouldNotBeNil)
				So(yes, ShouldBeFalse)
			})
		})
	})
}
