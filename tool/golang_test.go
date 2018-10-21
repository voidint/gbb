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
)

func TestBuild4Go(t *testing.T) {
	Convey("调用go build编译", t, func() {
		Convey("变量求值出错", func() {
			conf := config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Variable: "Hello", Value: "wrong express"},
				},
				Debug: true,
				All:   true,
			}
			setupConfig(&conf)
			builder := NewGoBuilder(conf)
			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldBeNil)
		})

		Convey("切换目录出错", func() {
			conf := config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "{{.Date}}"},
				},
				Debug: true,
				All:   true,
			}
			setupConfig(&conf)
			builder := NewGoBuilder(conf)

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
			monkey.Patch(util.WalkPkgsFunc, func(rootDir string, f util.FiltePkgFunc) (paths, symlinks []string, err error) {
				return nil, nil, ErrWalk
			})
			defer monkey.Unpatch(util.WalkPkgsFunc)

			conf := config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Debug:      true,
				All:        true,
			}
			setupConfig(&conf)
			builder := NewGoBuilder(conf)
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

			conf := config.Config{
				Tool:       "go build -ldflags='-w'",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []config.Variable{
					{Value: "{{.Date}}"},
				},
				Debug: true,
			}
			setupConfig(&conf)
			builder := NewGoBuilder(conf)

			wd, _ := os.Getwd()
			rootDir := filepath.Clean(strings.TrimSuffix(wd, "tool"))
			So(builder.Build(rootDir), ShouldBeNil)
		})
	})
}
