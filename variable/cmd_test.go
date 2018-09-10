package variable

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCmdVarMatch(t *testing.T) {
	Convey("命令形式的变量表达式匹配", t, func() {
		So(NewCmdVar().Match("$(date)"), ShouldBeTrue)
		So(NewCmdVar().Match("$(  date)"), ShouldBeTrue)
		So(NewCmdVar().Match("$($(date)  )"), ShouldBeTrue)

		So(NewCmdVar().Match("${date}"), ShouldBeFalse)
		So(NewCmdVar().Match("   $(date) "), ShouldBeFalse)
	})
}

func TestCmdVarEval(t *testing.T) {
	Convey("命令形式的变量表达式求值", t, func() {
		Convey("执行命令成功", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
				return []byte(strings.Join(c.Args, " ")), nil
			})
			defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

			val, err := NewCmdVar().Eval("", true)
			So(err, ShouldBeNil)
			So(val, ShouldBeBlank)

			if runtime.GOOS != "windows" {
				shell := os.Getenv("SHELL")
				defer os.Setenv("SHELL", shell)

				os.Setenv("SHELL", "/bin/bash")
				val, err = NewCmdVar().Eval("  $(date)", true)
				So(err, ShouldBeNil)
				So(val, ShouldEqual, fmt.Sprintf("/bin/bash -c %s", "date"))
			} else {
				val, err = NewCmdVar().Eval("  $(date)", true)
				So(err, ShouldBeNil)
				So(val, ShouldEqual, "date")
			}
		})

		Convey("执行命令异常", func() {
			var ErrExec = errors.New("exec cmd error")
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
				return nil, ErrExec
			})
			defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

			val, err := NewCmdVar().Eval("$(date)", true)
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, ErrExec)
			So(val, ShouldBeBlank)
		})
	})
}

func Test_exec(t *testing.T) {
	var cmd *exec.Cmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
		return []byte(strings.Join(c.Args, " ")), nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

	Convey("执行命令并获得执行的stdout", t, func() {
		if runtime.GOOS != "windows" {
			Convey("通过shell方式执行命令", func() {
				shell := os.Getenv("SHELL")

				os.Setenv("SHELL", "/bin/bash")
				defer os.Setenv("SHELL", shell)

				cmds := "git branch | grep '*' | awk {'print $2'}"
				out, err := NewCmdVar().exec(cmds, true)
				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, fmt.Sprintf("/bin/bash -c %s", cmds))
			})
		}
	})
}

func Test_execByNative(t *testing.T) {
	Convey("通过原生方式执行命令", t, func() {
		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
			return []byte(strings.Join(c.Args, " ")), nil
		})
		defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

		Convey("待执行的命令入参为空字符串", func() {
			nameAndArgs := ""
			So(func() {
				_, _ = NewCmdVar().execByNative(nameAndArgs, true)
			}, ShouldPanicWith, "unreachable")
		})
		Convey("待执行的命令仅包含命令名称", func() {
			nameAndArgs := "date"
			out, err := NewCmdVar().execByNative(nameAndArgs, true)
			So(err, ShouldBeNil)
			So(out, ShouldNotBeNil)
			So(string(out), ShouldEqual, "date")
		})
		Convey("待执行的命令含命令名称、选项和参数", func() {
			nameAndArgs := "go get -u -v github.com/voidint/gbb"
			out, err := NewCmdVar().execByNative(nameAndArgs, true)
			So(err, ShouldBeNil)
			So(out, ShouldNotBeNil)
			So(string(out), ShouldEqual, "go get -u -v github.com/voidint/gbb")
		})
	})
}

func Test_execByShell(t *testing.T) {
	Convey(`通过shell执行命令，如/bin/bash -c "echo 'hello world'"`, t, func() {
		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
			return []byte(strings.Join(c.Args, " ")), nil
		})
		defer monkey.UnpatchAll()

		sh, cmds := "/bin/bash", "git branch | grep '*' | awk {'print $2'}"
		out, err := NewCmdVar().execByShell(sh, cmds, true)
		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, fmt.Sprintf("%s -c %s", sh, cmds))
	})
}
