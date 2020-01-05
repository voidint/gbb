package tool

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
)

func TestBuild4GB(t *testing.T) {
	Convey("调用gb编译项目", t, func() {
		Convey("变量求值出错", func() {
			builder := NewGBBuilder(config.Config{
				Tool:       "gb build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "wrong express"},
				},
				Debug: true,
			})
			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldBeNil)
		})

		Convey("切换目录出错", func() {
			builder := NewGBBuilder(config.Config{
				Tool:       "gb build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "{{.Date}}"},
				},
				Debug: true,
			})

			var ErrChdir = errors.New("chdir error")
			monkey.Patch(util.Chdir, func(dir string, debug bool) (err error) {
				return ErrChdir
			})
			defer monkey.Unpatch(util.Chdir)

			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldEqual, ErrChdir)
		})

		Convey("编译成功", func() {
			builder := NewGBBuilder(config.Config{
				Tool:       "gb build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "{{.Date}}"},
				},
				Debug: true,
			})
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(_ *exec.Cmd) error {
				return nil
			})
			defer monkey.UnpatchAll()

			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldBeNil)
		})
	})
}
