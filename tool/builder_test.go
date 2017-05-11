package tool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/variable"
)

func initDir() {
	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, "tool") {
		return
	}
	idx := strings.Index(wd, filepath.Join("github.com", "voidint", "gbb"))
	if idx < 0 {
		return
	}
	wd = filepath.Join(wd[:idx], "github.com", "voidint", "gbb", "tool")
	_ = os.Chdir(wd)
}

func TestBuild(t *testing.T) {
	var cmd *exec.Cmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Run", func(c *exec.Cmd) error {
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(cmd), "Run")

	Convey("调用gb build编译", t, func() {
		initDir()
		wd, err := os.Getwd()
		So(err, ShouldBeNil)
		So(wd, ShouldNotBeBlank)
		So(strings.HasSuffix(wd, "tool"), ShouldBeTrue)

		Convey("包含变量表达式", func() {
			c := &config.Config{
				Version:    "0.3.0",
				Tool:       "gb build",
				Importpath: "github.com/voidint/gbb/build",
				Debug:      true,
			}

			Convey("包含非法变量表达式", func() {
				c.Variables = []config.Variable{
					{Variable: "Date", Value: "xxxx"},
				}
				err := Build(c, strings.TrimRight(wd, "tool"))
				So(err, ShouldNotBeNil)
			})

			Convey("包含合法变量表达式", func() {
				c.Variables = []config.Variable{
					{Variable: "Date", Value: "{{.Date}}"},
				}
				err := Build(c, "./")
				So(err, ShouldBeNil)
			})

		})

		Convey("不包含变量表达式", func() {
			c := &config.Config{
				Version: "0.3.0",
				Tool:    "gb build",
				Debug:   true,
			}

			err := Build(c, "./")
			So(err, ShouldBeNil)
		})
	})

	Convey("调用go build编译", t, func() {
		c := &config.Config{
			Version:    "0.3.0",
			Tool:       "go build",
			Importpath: "github.com/voidint/gbb/build",
			Variables: []config.Variable{
				{Variable: "Date", Value: "{{.Date}}"},
			},
			Debug: true,
		}

		err := Build(c, "./")
		So(err, ShouldBeNil)

	})

	Convey("调用非法的编译工具编译", t, func() {
		c := &config.Config{
			Version: "0.3.0",
			Tool:    "unsupported tool",
			Debug:   true,
		}

		err := Build(c, "./")
		So(err, ShouldEqual, ErrBuildTool)
	})
}

func TestChdir(t *testing.T) {
	Convey("切换工作目录", t, func() {
		Convey("目标目录是当前目录", func() {
			wd, err := os.Getwd()
			So(err, ShouldBeNil)
			So(chdir(wd, true), ShouldBeNil)
		})

		Convey("目标目录非当前目录", func() {
			wd, err := os.Getwd()
			So(err, ShouldBeNil)

			defer chdir(wd, true) // init work directory

			if idx := strings.LastIndex(wd, fmt.Sprintf("%c", os.PathSeparator)); idx > 0 {
				So(chdir(wd[:idx], true), ShouldBeNil)
			}
		})

		Convey("目录切换发生错误", func() {
			var ErrChdir = errors.New("chdir error")
			monkey.Patch(os.Getwd, func() (dir string, err error) {
				return "", ErrChdir
			})
			defer monkey.Unpatch(os.Getwd)
			So(chdir("../", true), ShouldEqual, ErrChdir)
		})
	})
}

func TestLdflags(t *testing.T) {
	Convey("根据配置返回-ldflags选项值", t, func() {
		Convey("从Tool中获取-ldflags选项值", func() {
			conf := new(config.Config)
			var flags string
			var err error

			conf.Tool = "go install"
			flags, err = ldflags(conf)
			So(err, ShouldBeNil)
			So(flags, ShouldBeEmpty)

			conf.Tool = "go install -ldflags='-w'"
			flags, err = ldflags(conf)
			So(err, ShouldBeNil)
			So(flags, ShouldEqual, "-w")
		})

		Convey("从变量中获取-ldflags选项值", func() {
			now := time.Unix(12345, 0)
			monkey.Patch(time.Now, func() time.Time {
				return now
			})

			var commitVar *variable.GitCommitVar
			hash := "abcdef12345"
			monkey.PatchInstanceMethod(reflect.TypeOf(commitVar), "Eval", func(_ *variable.GitCommitVar, _ string, debug bool) (val string, err error) {
				return hash, nil
			})
			defer monkey.UnpatchAll()

			conf := new(config.Config)
			var flags string
			var err error

			conf.Importpath = "github.com/voidint/gbb/build"
			conf.Variables = []config.Variable{
				{
					Variable: "Date",
					Value:    "{{.Date}}",
				},
				{
					Variable: "Commit",
					Value:    "{{.GitCommit}}",
				},
			}

			flags, err = ldflags(conf)
			So(err, ShouldBeNil)
			So(flags, ShouldEqual, fmt.Sprintf("-X %q -X %q",
				fmt.Sprintf("%s.Date=%s", conf.Importpath, now.Format(time.RFC3339)),
				fmt.Sprintf("%s.Commit=%s", conf.Importpath, hash),
			))

		})

	})
}

func TestExtractLdflags(t *testing.T) {
	Convey("抽取-ldflags选项值", t, func() {
		Convey("不包含该选项", func() {
			So(Args([]string{"go", "build"}).ExtractLdflags(), ShouldBeEmpty)
		})

		Convey("存在多个该选项及其值", func() {
			val := `-X "github.com/voidint/gbb/build.Date=2017-05-01T17:29:47+08:00" -X "github.com/voidint/gbb/build.Commit=07eba0ebf4648b9562182b682db12572da28f158"`
			So(Args([]string{
				"go",
				"build",
				fmt.Sprintf("-ldflags='%s'", val),
				"-ldflags",
				"'-w'",
			}).ExtractLdflags(), ShouldEqual, val)
		})

		Convey("存在一个该选项及其值", func() {
			Convey("选项与值之间使用'='符号分隔", func() {
				So(Args([]string{"go", "build", "-ldflags='-w'"}).ExtractLdflags(), ShouldEqual, "-w")
			})
			Convey("选项与值之间使用空格分隔", func() {
				So(Args([]string{"go", "build", "-ldflags", "\"-w\""}).ExtractLdflags(), ShouldEqual, "-w")
				So(Args([]string{"go", "build", "-ldflags"}).ExtractLdflags(), ShouldBeEmpty)
			})
			Convey("值放置在两个单引号内", func() {
				So(Args([]string{"go", "build", "-ldflags", "'-w'"}).ExtractLdflags(), ShouldEqual, "-w")
			})
			Convey("值放置在两个双引号内", func() {
				So(Args([]string{"go", "build", "-ldflags=\"-w\""}).ExtractLdflags(), ShouldEqual, "-w")
			})
		})
	})
}

func TestRemoveLdflags(t *testing.T) {
	Convey("移除参数中的-ldflags", t, func() {
		var news []string
		news = Args([]string{"go", "build"}).RemoveLdflags()
		So(len(news), ShouldEqual, 2)
		So(news[0], ShouldEqual, "go")
		So(news[1], ShouldEqual, "build")

		news = Args([]string{"go", "build", "-ldflags='-w'", "-gcflags='-N -l'"}).RemoveLdflags()
		So(len(news), ShouldEqual, 3)
		So(news[0], ShouldEqual, "go")
		So(news[1], ShouldEqual, "build")
		So(news[2], ShouldEqual, "-gcflags='-N -l'")

		news = Args([]string{"go", "build", "-ldflags", "'-w'", "-gcflags='-N -l'"}).RemoveLdflags()
		So(len(news), ShouldEqual, 3)
		So(news[0], ShouldEqual, "go")
		So(news[1], ShouldEqual, "build")
		So(news[2], ShouldEqual, "-gcflags='-N -l'")
	})
}

func TestTrimQuotationMarks(t *testing.T) {
	Convey("去除前后单/双引号", t, func() {
		So(TrimQuotationMarks(`'-w'`), ShouldEqual, "-w")
		So(TrimQuotationMarks(`"-N -l"`), ShouldEqual, "-N -l")
		So(TrimQuotationMarks(`something`), ShouldEqual, "something")

		So(TrimQuotationMarks(`'-w"`), ShouldEqual, "'-w")
		So(TrimQuotationMarks(`"-w'`), ShouldEqual, "\"-w")
	})
}
