package variable

import (
	"errors"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/lmika/shellwords"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewGitCommitVar(t *testing.T) {
	Convey("创建GitCommit变量实例", t, func() {
		v, err := NewGitCommitVar(GitCmd, "")
		So(err, ShouldEqual, ErrExprEmpty)
		So(v, ShouldBeNil)
		v, err = NewGitCommitVar("", DefaultGitCommitExpr)
		So(err, ShouldEqual, ErrCmdEmpty)
		So(v, ShouldBeNil)
		v, err = NewGitCommitVar(GitCmd, DefaultGitCommitExpr)
		So(err, ShouldBeNil)
		So(v, ShouldResemble, DefaultGitCommitVar)
	})
}

func TestGitCommitVarMatch(t *testing.T) {
	Convey("内置GitCommit变量表达式匹配", t, func() {
		So(DefaultGitCommitVar.Match(DefaultGitCommitExpr), ShouldBeTrue)
		v, err := NewGitCommitVar(GitCmd, "anything")
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)

		So(v.Match("anything"), ShouldBeTrue)
		So(v.Match("$(  date)"), ShouldBeFalse)
	})
}

func TestGitCommitVarEval(t *testing.T) {
	Convey("内置GitCommit变量表达式求值", t, func() {
		Convey("求值成功", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
				return []byte(strings.Join(c.Args, " ")), nil
			})
			defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

			val, err := DefaultGitCommitVar.Eval("", true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, GitCmd)

			v, _ := NewGitCommitVar("anything", DefaultGitCommitExpr)
			val, err = v.Eval("", true)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "anything")
		})

		Convey("求值失败", func() {
			var ErrEval = errors.New("eval error")
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(c *exec.Cmd) ([]byte, error) {
				return nil, ErrEval
			})
			defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Output")

			val, err := DefaultGitCommitVar.Eval("", true)
			So(err, ShouldEqual, ErrEval)
			So(val, ShouldBeBlank)
		})

		Convey("panic", func() {
			monkey.Patch(shellwords.Split, func(s string) []string {
				return []string{}
			})
			defer monkey.Unpatch(shellwords.Split)
			So(func() {
				_, _ = DefaultGitCommitVar.Eval("", true)
			}, ShouldPanicWith, "unreachable")
		})
	})
}
