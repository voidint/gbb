package tool

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
	"github.com/voidint/gbb/variable"
)

func TestBuild4Go(t *testing.T) {
	Convey("调用go build编译", t, func() {
		Convey("变量求值出错", func() {
			builder := NewGoBuilder(config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "wrong express"},
				},
				Debug: true,
				All:   true,
			})
			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldEqual, variable.ErrExpr)
		})

		Convey("切换目录出错", func() {
			builder := NewGoBuilder(config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "{{.Date}}"},
				},
				Debug: true,
				All:   true,
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

		Convey("遍历所有go package出错", func() {
			var ErrWalk = errors.New("walk error")
			monkey.Patch(util.WalkPkgsFunc, func(rootDir string, f util.FiltePkgFunc) (paths []string, err error) {
				return nil, ErrWalk
			})
			defer monkey.Unpatch(util.WalkPkgsFunc)

			builder := NewGoBuilder(config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Debug:      true,
				All:        true,
			})
			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldEqual, ErrWalk)
		})

		Convey("遍历所有go package成功并执行go clean", func() {
			var cmd *exec.Cmd
			monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(_ *exec.Cmd) error {
				return nil
			})
			defer monkey.UnpatchAll()

			builder := NewGoBuilder(config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "{{.Date}}"},
				},
				Debug: true,
			})

			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldBeNil)
		})
	})
}
