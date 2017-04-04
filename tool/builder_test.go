package tool

import (
	"fmt"
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

func initDir() {
	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, "tool") {
		return
	}
	idx := strings.Index(wd, filepath.Join("github.com", "voidint", "gbb"))
	if idx < 0 {
		return
	}
	wd = filepath.Join(wd[:idx], "github.com", "voidint", "gbb", "tool")
	_ = os.Chdir(wd)
}

func TestBuild(t *testing.T) {
	var cmd *exec.Cmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(c *exec.Cmd) error {
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Run")

	Convey("调用gb build编译", t, func() {
		initDir()
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(wd, ShouldNotBeBlank)
		So(strings.HasSuffix(wd, "tool"), ShouldBeTrue)

		Convey("包含变量表达式", func() {
			c := &config.Config{
				Version:    "0.3.0",
				Tool:       "gb build",
				Importpath: "github.com/voidint/gbb/build",
				Debug:      true,
			}

			Convey("包含非法变量表达式", func() {
				c.Variables = []config.Variable{
					{Variable: "Date", Value: "xxxx"},
				}
				err := Build(c, strings.TrimRight(wd, "tool"))
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
			Version:    "0.3.0",
			Tool:       "go build",
			Importpath: "github.com/voidint/gbb/build",
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

func TestChdir(t *testing.T) {
	Convey("切换工作目录", t, func() {
		Convey("目标目录是当前目录", func() {
			wd, err := os.Getwd()
			So(err, ShouldBeNil)
			So(chdir(wd, true), ShouldBeNil)
		})

		Convey("目标目录非当前目录", func() {
			wd, err := os.Getwd()
			So(err, ShouldBeNil)

			defer chdir(wd, true) // init work directory

			if idx := strings.LastIndex(wd, fmt.Sprintf("%c", os.PathSeparator)); idx > 0 {
				So(chdir(wd[:idx], true), ShouldBeNil)
			}
		})
	})
}
