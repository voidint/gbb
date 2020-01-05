package variable

import (
	"testing"
	"time"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func now() time.Time {
	return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
}

func TestDateVarMatch(t *testing.T) {
	Convey("内置Date变量表达式匹配", t, func() {
		So(NewDateVar(time.RFC3339, DefaultDateExpr).Match(DefaultDateExpr), ShouldBeTrue)
		So(NewDateVar(time.RFC3339, DefaultDateExpr).Match("$(  date)"), ShouldBeFalse)
	})
}

func TestDateVarEval(t *testing.T) {
	Convey("内置Date变量表达式求值", t, func() {
		monkey.Patch(time.Now, now)
		defer monkey.Unpatch(time.Now)

		val, err := NewDateVar("", "").Eval("", true)
		So(err, ShouldBeNil)
		So(val, ShouldEqual, "")

		val, err = NewDateVar(time.RFC3339, "{{.Date}}").Eval("", true)
		So(err, ShouldBeNil)
		So(val, ShouldEqual, now().Format(time.RFC3339))
	})
}
