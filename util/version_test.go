package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersionGreaterThan(t *testing.T) {
	Convey("判断两个语义化版本号的大小关", t, func() {
		var yes bool
		var err error
		Convey("参数中包含非语义化版本号", func() {
			yes, err = VersionGreaterThan("1.0.5.6", "0.0.1")
			So(err, ShouldEqual, ErrSemanticVersion)
			So(yes, ShouldBeFalse)

			yes, err = VersionGreaterThan("0.1.2", "v0.0.1")
			So(err, ShouldEqual, ErrSemanticVersion)
			So(yes, ShouldBeFalse)

			yes, err = VersionGreaterThan("a.b.c", "0.0.1")
			So(err, ShouldEqual, ErrSemanticVersion)
			So(yes, ShouldBeFalse)
		})

		Convey("参数都满足语义化版本号要求", func() {
			yes, err = VersionGreaterThan("8.5.6", "9.10.1")
			So(err, ShouldBeNil)
			So(yes, ShouldBeFalse)

			yes, err = VersionGreaterThan("10.5.6", "9.10.1")
			So(err, ShouldBeNil)
			So(yes, ShouldBeTrue)

			yes, err = VersionGreaterThan("0.5.6", "0.10.1")
			So(err, ShouldBeNil)
			So(yes, ShouldBeFalse)

			yes, err = VersionGreaterThan("0.5.6", "0.1.1")
			So(err, ShouldBeNil)
			So(yes, ShouldBeTrue)

			yes, err = VersionGreaterThan("0.5.6", "0.5.5")
			So(err, ShouldBeNil)
			So(yes, ShouldBeTrue)

			yes, err = VersionGreaterThan("0.5.6", "0.5.6")
			So(err, ShouldBeNil)
			So(yes, ShouldBeFalse)
		})

	})
}
