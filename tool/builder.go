package tool

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/variable"
)

// Builder 编译工具
type Builder interface {
	Build(dir string) error
}

var (
	// ErrBuildTool 不支持的编译工具错误
	ErrBuildTool = errors.New("unsupported build tool")
)

// Build 根据配置信息，调用合适的编译工具进行编译。
// 若配置的编译工具不在支持的工具范围内，则返回ErrBuildTool错误。
func Build(conf *config.Config, debug bool, dir string) error {
	if strings.HasPrefix(conf.Tool, "go ") {
		return NewGoBuilder(conf, debug).Build(dir)
	} else if strings.HasPrefix(conf.Tool, "gb ") {
		return NewGBBuilder(conf, debug).Build(dir)
	}
	return ErrBuildTool
}

// buildDir 切换到指定工作目录，调用指定的编译工具进行编译。
func buildDir(conf *config.Config, debug bool, dir string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if wd != dir {
		if debug {
			fmt.Printf("==> cd %s\n", dir)
		}
		if err = os.Chdir(dir); err != nil {
			return err
		}
	}

	cmdArgs := strings.Fields(conf.Tool) // go install ==> []string{"go", "install"}

	var buf bytes.Buffer
	for i := range conf.Variables {
		varName := conf.Variables[i].Variable
		varExpr := conf.Variables[i].Value

		val, err := variable.Eval(varExpr)
		if err != nil {
			return err
		}
		buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Package, varName, val))
		if i < len(conf.Variables)-1 {
			buf.WriteByte(' ')
		}
	}
	if buf.Len() > 0 {
		cmdArgs = append(cmdArgs, "-ldflags", buf.String())
	}

	if debug {
		fmt.Print("==> ", cmdArgs[0])
		args := cmdArgs[1:]
		for i := range args {
			if i-1 > 0 && args[i-1] == "-ldflags" {
				fmt.Printf(" '%s'", args[i])
			} else {
				fmt.Printf(" %s", args[i])
			}
		}
		fmt.Println()
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
