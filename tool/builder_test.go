package tool

import (
	"os/exec"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
)

func TestBuild(t *testing.T) {
	var cmd *exec.Cmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(c *exec.Cmd) error {
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Run")

	Convey("调用gb build编译", t, func() {
		Convey("包含变量表达式", func() {
			c := &config.Config{
				Version: "0.3.0",
				Tool:    "gb build",
				Package: "github.com/voidint/gbb/build",
				Debug:   true,
			}

			Convey("包含非法变量表达式", func() {
				c.Variables = []config.Variable{
					{Variable: "Date", Value: "xxxx"},
				}
				err := Build(c, "./")
				So(err, ShouldNotBeNil)
			})

			Convey("包含合法变量表达式", func() {
				c.Variables = []config.Variable{
					{Variable: "Date", Value: "{{.Date}}"},
				}
				err := Build(c, "./")
				So(err, ShouldBeNil)
			})

		})

		Convey("不包含变量表达式", func() {
			c := &config.Config{
				Version: "0.3.0",
				Tool:    "gb build",
				Debug:   true,
			}

			err := Build(c, "./")
			So(err, ShouldBeNil)
		})
	})

	Convey("调用go build编译", t, func() {
		c := &config.Config{
			Version: "0.3.0",
			Tool:    "go build",
			Package: "github.com/voidint/gbb/build",
			Variables: []config.Variable{
				{Variable: "Date", Value: "{{.Date}}"},
			},
			Debug: true,
		}

		err := Build(c, "./")
		So(err, ShouldBeNil)

	})

	Convey("调用非法的编译工具编译", t, func() {
		c := &config.Config{
			Version: "0.3.0",
			Tool:    "unsupported tool",
			Debug:   true,
		}

		err := Build(c, "./")
		So(err, ShouldEqual, ErrBuildTool)
	})
}
