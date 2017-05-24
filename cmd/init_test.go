package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
)

func TestGatherOne(t *testing.T) {
	Convey("收集用户终端输入", t, func() {
		Convey("带默认值，预期输入go_build", func() {
			monkey.Patch(fmt.Scanln, func(a ...interface{}) (n int, err error) {
				val := "go_build"
				input := a[0].(*string)
				*input = val
				return len(val), nil
			})
			defer monkey.Unpatch(fmt.Scanln)

			So(gatherOne("importpath", "build"), ShouldEqual, "go build")
		})

		Convey("带默认值，预期输入为空", func() {
			monkey.Patch(fmt.Scanln, func(a ...interface{}) (n int, err error) {
				return 0, nil
			})
			defer monkey.Unpatch(fmt.Scanln)

			So(gatherOne("tool", "go_install"), ShouldEqual, "go install")
		})

		Convey("不带默认值，预期输入github.com/voidint/gbb/build", func() {
			monkey.Patch(fmt.Scanln, func(a ...interface{}) (n int, err error) {
				val := "github.com/voidint/gbb/build"
				input := a[0].(*string)
				*input = val
				return len(val), nil
			})
			defer monkey.Unpatch(fmt.Scanln)

			So(gatherOne("tool", ""), ShouldEqual, "github.com/voidint/gbb/build")
		})

	})
}

func TestGatherOneVar(t *testing.T) {
	Convey("收集用户输入的变量名及其值", t, func() {
		monkey.Patch(fmt.Scanln, func(a ...interface{}) (n int, err error) {
			val := "something"
			input := a[0].(*string)
			*input = val
			return len(val), nil
		})
		defer monkey.Unpatch(fmt.Scanln)

		v := gatherOneVar()
		So(v, ShouldNotBeNil)
		So(v.Variable, ShouldEqual, "something")
		So(v.Value, ShouldEqual, "something")
	})
}

func TestGather(t *testing.T) {
	Convey("收集用户的多次输入", t, func() {
		monkey.Patch(fmt.Scanln, func(a ...interface{}) (n int, err error) {
			val := "n"
			input := a[0].(*string)
			*input = val
			return len(val), nil
		})
		defer monkey.Unpatch(fmt.Scanln)

		c := gather()
		So(c, ShouldNotBeNil)
		So(c.Version, ShouldEqual, Version)
		So(c.Tool, ShouldEqual, "n")
	})
}

func TestGenConfigFile(t *testing.T) {
	Convey("在指定路径生成配置文件", t, func() {
		monkey.Patch(gather, func() (c *config.Config) {
			return &config.Config{
				Version:    Version,
				Tool:       "go build",
				Importpath: "github.com/voidint/gbb/build",
			}
		})

		monkey.Patch(fmt.Scanln, func(a ...interface{}) (n int, err error) {
			val := "y"
			input := a[0].(*string)
			*input = val
			return len(val), nil
		})

		defer monkey.UnpatchAll()

		genConfigFile("./gbb_test.json")
		os.Remove("./gbb_test.json")
	})
}

func TestInitCmd(t *testing.T) {
	monkey.Patch(genConfigFile, func(_ string) {
	})

	defer monkey.Unpatch(genConfigFile)
	initCmd.Run(nil, nil)
}
