package cmd

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/tool"
	"github.com/voidint/gbb/util"
)

func Test_initConfig(t *testing.T) {
	Convey("全局变量初始化", t, func() {
		monkey.Patch(os.Getwd, func() (dir string, err error) {
			return "/a/b/c", nil
		})
		defer monkey.UnpatchAll()

		initConfig()
		So(wd, ShouldEqual, "/a/b/c")

		var ErrGetwd = errors.New("can't get working directory")
		monkey.Patch(os.Getwd, func() (dir string, err error) {
			return "", ErrGetwd
		})

		monkey.Patch(os.Exit, func(_ int) {
		})

		initConfig()
		So(wd, ShouldBeEmpty)
	})
}

func TestExecute(t *testing.T) {
	var ErrExec = errors.New("exec error")
	var cmd *cobra.Command
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Execute", func(_ *cobra.Command) error {
		return ErrExec
	})
	monkey.Patch(os.Exit, func(_ int) {
	})

	defer monkey.UnpatchAll()
	Execute()
}

func TestRootCmd(t *testing.T) {
	Convey("RootCmd.Run", t, func() {
		Convey("gbb.json不存在", func() {
			monkey.Patch(util.FileExist, func(filename string) bool {
				return false
			})
			monkey.Patch(genConfigFile, func(_ string) error {
				return nil
			})
			defer monkey.UnpatchAll()

			RootCmd.Run(nil, nil)
		})

		Convey("gbb.json已存在", func() {
			Convey("解析gbb.json文件出错", func() {
				monkey.Patch(util.FileExist, func(filename string) bool {
					return true
				})
				defer monkey.UnpatchAll()

				monkey.Patch(config.Load, func(_ string) (*config.Config, error) {
					return nil, errors.New("paser error")
				})
				monkey.Patch(os.Exit, func(code int) {
				})

				RootCmd.Run(nil, nil)
			})

			Convey("版本号比较出错", func() {
				monkey.Patch(util.FileExist, func(filename string) bool {
					return true
				})
				defer monkey.UnpatchAll()
				monkey.Patch(util.VersionGreaterThan, func(_, _ string) (bool, error) {
					return false, errors.New("version error")
				})
				monkey.Patch(config.Load, func(_ string) (conf *config.Config, err error) {
					return &config.Config{
						Version: "0.0.1",
					}, nil
				})
				monkey.Patch(os.Exit, func(code int) {
				})
				RootCmd.Run(nil, nil)
			})
			Convey("gbb.json版本号高于gbb程序版本号", func() {
				monkey.Patch(util.FileExist, func(filename string) bool {
					return true
				})
				monkey.Patch(config.Load, func(_ string) (conf *config.Config, err error) {
					return &config.Config{
						Version: "999.999.999",
					}, nil
				})
				monkey.Patch(tool.Build, func(_ *config.Config, _ string) (err error) {
					return nil
				})
				defer monkey.UnpatchAll()
				RootCmd.Run(nil, nil)
			})

			Convey("gbb.json版本号低于gbb程序版本号", func() {
				monkey.Patch(util.FileExist, func(filename string) bool {
					return true
				})
				monkey.Patch(config.Load, func(_ string) (conf *config.Config, err error) {
					return &config.Config{
						Version: "0.0.0",
					}, nil
				})
				monkey.Patch(genConfigFile, func(_ string) error {
					return nil
				})
				monkey.Patch(tool.Build, func(_ *config.Config, _ string) (err error) {
					return nil
				})
				defer monkey.UnpatchAll()
				RootCmd.Run(nil, nil)
			})
		})

		Convey("编译失败", func() {
			monkey.Patch(util.FileExist, func(filename string) bool {
				return true
			})
			monkey.Patch(config.Load, func(_ string) (conf *config.Config, err error) {
				return &config.Config{
					Version: Version,
				}, nil
			})
			monkey.Patch(tool.Build, func(_ *config.Config, _ string) (err error) {
				return errors.New("build error")
			})
			monkey.Patch(os.Exit, func(_ int) {
			})
			defer monkey.UnpatchAll()
			RootCmd.Run(nil, nil)
		})
	})
}
