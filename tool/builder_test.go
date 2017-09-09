package tool

import (
	"bytes"
	"errors"
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
	"github.com/voidint/gbb/variable"
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
				err := Build(c, strings.TrimSuffix(wd, "tool"))
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

func Test_setupVars(t *testing.T) {
	const (
		dateName   = "Date"
		dateVal    = "2017-09-03T16:58:19+08:00"
		commitName = "GitCommit"
		commitVal  = "9a6869e1591752d29973535b48a5ecfe7471eb49"
	)

	Convey("设置config指针中的变量值", t, func() {
		Convey("未定义变量，无需对变量求值", func() {
			So(setupVars(&config.Config{}), ShouldBeNil)
		})

		Convey("变量求值成功", func() {
			monkey.Patch(variable.Eval, func(expr string, debug bool) (val string, err error) {
				switch expr {
				case variable.DefaultDateExpr:
					return dateVal, nil
				case variable.DefaultGitCommitExpr:
					return commitVal, nil
				}
				panic("unreachable")
			})
			defer monkey.Unpatch(variable.Eval)
			conf := config.Config{
				Variables: []config.Variable{
					{Variable: dateName, Value: variable.DefaultDateExpr},
					{Variable: commitName, Value: variable.DefaultGitCommitExpr},
				},
			}
			So(setupVars(&conf), ShouldBeNil)
			So(conf.Variables[0].Value, ShouldEqual, dateVal)
			So(conf.Variables[1].Value, ShouldEqual, commitVal)
		})

		Convey("变量求值发生错误", func() {
			ErrEval := errors.New("eval error")
			monkey.Patch(variable.Eval, func(expr string, debug bool) (val string, err error) {
				return "", ErrEval
			})
			defer monkey.Unpatch(variable.Eval)

			conf := config.Config{
				Variables: []config.Variable{
					{Variable: dateName, Value: variable.DefaultDateExpr},
					{Variable: commitName, Value: variable.DefaultGitCommitExpr},
				},
			}
			So(setupVars(&conf), ShouldEqual, ErrEval)
		})
	})
}

func Test_setupTool(t *testing.T) {
	const (
		dateName   = "Date"
		dateVal    = "2017-09-03T16:58:19+08:00"
		commitName = "GitCommit"
		commitVal  = "9a6869e1591752d29973535b48a5ecfe7471eb49"
	)
	Convey("设置config指针中的tool值", t, func() {
		// Convey("未定义变量", func() {
		// 	setupTool(&config.Config{})
		// })

		Convey("未设置-ldflags选项", func() {
			conf := config.Config{
				Tool:       "go build",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Variable: dateName, Value: dateVal},
					{Variable: commitName, Value: commitVal},
				},
			}

			var buf bytes.Buffer
			for i := range conf.Variables {
				buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Importpath, conf.Variables[i].Variable, conf.Variables[i].Value))
				if i < len(conf.Variables)-1 {
					buf.WriteByte(' ')
				}
			}

			setupTool(&conf)
			So(conf.Tool, ShouldEqual, fmt.Sprintf("go build -ldflags '%s'", buf.String()))
		})

		Convey("已设置-ldflags选项", func() {
			conf := config.Config{
				Tool:       "go build -ldflags='-s -w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Variable: dateName, Value: dateVal},
					{Variable: commitName, Value: commitVal},
				},
			}

			var buf bytes.Buffer
			for i := range conf.Variables {
				buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Importpath, conf.Variables[i].Variable, conf.Variables[i].Value))
				if i < len(conf.Variables)-1 {
					buf.WriteByte(' ')
				}
			}

			setupTool(&conf)
			So(conf.Tool, ShouldEqual, fmt.Sprintf("go build -ldflags '-s -w %s'", buf.String()))
		})
	})
}
