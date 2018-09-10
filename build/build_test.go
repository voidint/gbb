package build

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	Convey("生成版本信息", t, func() {
		prefix := "gbb version 0.3.0"
		Convey("Date、Commit变量都为空", func() {
			So(Version(prefix), ShouldEqual, prefix)
		})
		Convey("Date、Branch、Commit变量全非空", func() {
			Date = "2017-01-09T09:14:54+08:00"
			Branch = "master"
			Commit = "ec7dd797369606461ac9e861e771b730321f3e2f"
			val := fmt.Sprintf("gbb version 0.3.0\ndate: %s\nbranch: %s\ncommit: %s", Date, Branch, Commit)
			So(Version(prefix), ShouldEqual, val)
		})
		Convey("Date、Commit变量非全空", func() {
			Convey("Date为空，Branch、Commit非空", func() {
				Date = ""
				Branch = "master"
				Commit = "ec7dd797369606461ac9e861e771b730321f3e2f"
				val := fmt.Sprintf("gbb version 0.3.0\nbranch: %s\ncommit: %s", Branch, Commit)
				So(Version(prefix), ShouldEqual, val)
			})
			Convey("Branch为空，Date、Commit非空", func() {
				Branch = ""
				Date = "2017-01-09T09:14:54+08:00"
				Commit = "ec7dd797369606461ac9e861e771b730321f3e2f"
				val := fmt.Sprintf("gbb version 0.3.0\ndate: %s\ncommit: %s", Date, Commit)
				So(Version(prefix), ShouldEqual, val)
			})
			Convey("Commit为空，Date、Branch非空", func() {
				Commit = ""
				Date = "2017-01-09T09:14:54+08:00"
				Branch = "master"
				val := fmt.Sprintf("gbb version 0.3.0\ndate: %s\nbranch: %s", Date, Branch)
				So(Version(prefix), ShouldEqual, val)
			})
		})
	})
}
