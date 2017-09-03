package tool

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
)

func TestBuildDir4Golang(t *testing.T) {
	builder := NewGoBuilder(config.Config{
		Tool:       "go build -ldflags='-w'",
		Importpath: "github.com/voidint/gbb/build",
		Debug:      true,
	})

	Convey("编译指定目录及其子目录下的go源文件", t, func() {
		dir, err := os.Getwd()
		So(err, ShouldBeNil)
		So(dir, ShouldNotBeEmpty)

		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(_ *exec.Cmd) error {
			return nil
		})
		defer monkey.UnpatchAll()

		So(builder.buildDir(dir), ShouldBeNil)
	})

	Convey("编译指定目录及其子目录下的go源文件出错", t, func() {
		Convey("切换目录出错", func() {
			var ErrChdir = errors.New("chdir error")
			monkey.Patch(chdir, func(dir string, debug bool) (err error) {
				return ErrChdir
			})
			defer monkey.Unpatch(chdir)
			So(builder.buildDir("../"), ShouldEqual, ErrChdir)
		})
	})
}

func TestWalkPkgs(t *testing.T) {
	Convey("查找指定目录及其子目录下所有满足golang包的目录路径", t, func() {
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(strings.HasSuffix(wd, "tool"), ShouldBeTrue)

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
		So(strings.HasSuffix(wd, "tool"), ShouldBeTrue)

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
