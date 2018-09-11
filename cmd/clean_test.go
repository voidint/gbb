package cmd

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
	"github.com/voidint/gbb/util"
)

func Test_clean(t *testing.T) {
	Convey("切换到指定目录后执行go clean", t, func() {
		Convey("切换工作目录出错", func() {
			var ErrChdir = errors.New("chdir error")
			monkey.Patch(util.Chdir, func(dir string, debug bool) (err error) {
				return ErrChdir
			})
			defer monkey.Unpatch(util.Chdir)

			beforeDir, _ := os.Getwd()
			So(clean("../", &CleanOptions{}, &GlobalOptions{}), ShouldEqual, ErrChdir)
			afterDir, _ := os.Getwd()
			So(beforeDir == afterDir, ShouldBeTrue)
		})

		Convey("成功切换到指定工作目录", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(_ *exec.Cmd) error {
				return nil
			})
			defer monkey.UnpatchAll()

			beforeDir, _ := os.Getwd()
			pwd := filepath.Clean(strings.TrimSuffix(beforeDir, "cmd"))
			So(clean(pwd, &CleanOptions{
				I: true,
				N: true,
				R: true,
				X: true,
			}, &GlobalOptions{Debug: true}), ShouldBeNil)
			afterDir, _ := os.Getwd()
			So(afterDir, ShouldEqual, pwd)
		})
	})
}

func Test_rclean(t *testing.T) {
	Convey("在指定目录及其子目录下寻找main package目录并在其中执行go clean", t, func() {
		Convey("遍历目录出错", func() {
			var ErrWalk = errors.New("walk error")
			monkey.Patch(util.WalkPkgsFunc, func(rootDir string, f util.FiltePkgFunc) (paths []string, err error) {
				return nil, ErrWalk
			})
			defer monkey.Unpatch(util.WalkPkgsFunc)
			beforeDir, _ := os.Getwd()
			So(rclean(beforeDir, &CleanOptions{}, &GlobalOptions{}), ShouldEqual, ErrWalk)
			afterDir, _ := os.Getwd()
			So(afterDir, ShouldEqual, beforeDir)
		})

		Convey("遍历目录及其子目录成功并执行go clean", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(_ *exec.Cmd) error {
				return nil
			})
			defer monkey.UnpatchAll()

			beforeDir, _ := os.Getwd()
			pwd := filepath.Clean(strings.TrimSuffix(beforeDir, "cmd"))
			So(rclean(pwd, &CleanOptions{}, &GlobalOptions{}), ShouldBeNil)
			afterDir, _ := os.Getwd()
			So(afterDir, ShouldEqual, pwd)
		})
	})
}

func TestCleanCmd(t *testing.T) {
	var ErrRclean = errors.New("rclean error")
	monkey.Patch(rclean, func(rootDir string, opts *CleanOptions, gopts *GlobalOptions) (err error) {
		return ErrRclean
	})

	monkey.Patch(os.Exit, func(code int) {
	})

	defer monkey.UnpatchAll()
	cleanCmd.Run(nil, nil)
}
