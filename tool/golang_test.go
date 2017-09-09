package tool

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
)

func TestBuildDir4Golang(t *testing.T) {
	builder := NewGoBuilder(config.Config{
		Tool:       "go build -ldflags='-w'",
		Importpath: "github.com/voidint/gbb/build",
		Debug:      true,
	})

	Convey("编译指定目录及其子目录下的go源文件", t, func() {
		dir, err := os.Getwd()
		So(err, ShouldBeNil)
		So(dir, ShouldNotBeEmpty)

		var cmd *exec.Cmd
		monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(_ *exec.Cmd) error {
			return nil
		})
		defer monkey.UnpatchAll()

		So(builder.buildDir(dir), ShouldBeNil)
	})

	Convey("编译指定目录及其子目录下的go源文件出错", t, func() {
		Convey("切换目录出错", func() {
			var ErrChdir = errors.New("chdir error")
			monkey.Patch(util.Chdir, func(dir string, debug bool) (err error) {
				return ErrChdir
			})
			defer monkey.Unpatch(util.Chdir)
			So(builder.buildDir("../"), ShouldEqual, ErrChdir)
		})
	})
}
