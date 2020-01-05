package variable

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEval(t *testing.T) {
	Convey("遍历所有内建的变量表达式列表求值", t, func() {
		Convey("匹配到Date变量表达式", func() {
			monkey.Patch(time.Now, now)
			defer monkey.Unpatch(time.Now)

			val, err := Eval(DefaultDateExpr, true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, now().Format(time.RFC3339))
		})

		Convey("匹配到GitCommit变量表达式", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
				return []byte(strings.Join(c.Args, " ")), nil
			})
			defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")
			val, err := Eval(DefaultGitCommitExpr, true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, GitCmd)
		})

		Convey("匹配到Cmd变量表达式", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
				return []byte(strings.Join(c.Args, " ")), nil
			})
			defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

			if runtime.GOOS != "windows" {
				shell := os.Getenv("SHELL")
				defer os.Setenv("SHELL", shell)

				os.Setenv("SHELL", "/bin/bash")
				val, err := NewCmdVar().Eval("$(date)", true)
				So(err, ShouldBeNil)
				So(val, ShouldEqual, fmt.Sprintf("/bin/bash -c %s", "date"))
			} else {
				val, err := NewCmdVar().Eval("$(date)", true)
				So(err, ShouldBeNil)
				So(val, ShouldEqual, "date")
			}
		})

		Convey("未匹配到任何变量表达式", func() {
			expr := "xxxxx"
			val, err := Eval(expr, true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, expr)
		})
	})
}
