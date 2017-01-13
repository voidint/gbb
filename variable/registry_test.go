package variable

import (
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bouk/monkey"
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

			val, err := Eval("$(Date)", true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "Date")
		})

		Convey("未匹配到任何变量表达式", func() {
			val, err := Eval("xxxxx", true)
			So(err, ShouldEqual, ErrExpr)
			So(val, ShouldEqual, "")
		})
	})
}
