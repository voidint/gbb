package variable

import (
	"errors"
	"os/exec"
	"reflect"
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

			val, err = NewCmdVar().Eval("  $(date)", true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "date")
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

func TestExec(t *testing.T) {
	var cmd *exec.Cmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
		return []byte(strings.Join(c.Args, " ")), nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")
	Convey("执行命令并获得执行的stdout", t, func() {
		Convey("待执行的命令入参为空字符串", func() {
			nameAndArgs := ""
			So(func() {
				_, _ = NewCmdVar().exec(nameAndArgs)
			}, ShouldPanicWith, "unreachable")
		})
		Convey("待执行的命令仅包含命令名称", func() {
			nameAndArgs := "date"
			out, err := NewCmdVar().exec(nameAndArgs)
			So(err, ShouldBeNil)
			So(out, ShouldNotBeNil)
			So(string(out), ShouldEqual, "date")
		})
		Convey("待执行的命令含命令名称、选项和参数", func() {
			nameAndArgs := "go get -u -v github.com/voidint/gbb"
			out, err := NewCmdVar().exec(nameAndArgs)
			So(err, ShouldBeNil)
			So(out, ShouldNotBeNil)
			So(string(out), ShouldEqual, "go get -u -v github.com/voidint/gbb")
		})
	})
}
