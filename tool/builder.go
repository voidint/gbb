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
func Build(conf *config.Config, dir string) (err error) {
	defer chdir(dir, conf.Debug) // init work directory

	if strings.HasPrefix(conf.Tool, "go ") {
		return NewGoBuilder(conf).Build(dir)
	} else if strings.HasPrefix(conf.Tool, "gb ") {
		return NewGBBuilder(conf).Build(dir)
	}
	return ErrBuildTool
}

// buildDir 切换到指定工作目录，调用指定的编译工具进行编译。
func buildDir(conf *config.Config, dir string) (err error) {
	if err = chdir(dir, conf.Debug); err != nil {
		return err
	}

	cmdArgs := strings.Fields(conf.Tool) // go install ==> []string{"go", "install"}

	mainPkg, err := isMainPkg(dir)
	if err != nil {
		return err
	}
	if mainPkg {
		flags, err := ldflags(conf)
		if err != nil {
			return err
		}
		if flags != "" {
			cmdArgs = append(cmdArgs, "-ldflags", flags)
		}
	}

	if conf.Debug {
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

func chdir(dir string, debug bool) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if wd == dir {
		return nil
	}

	if debug {
		fmt.Printf("==> cd %s\n", dir)
	}
	return os.Chdir(dir)
}

func ldflags(conf *config.Config) (flags string, err error) {
	var buf bytes.Buffer
	for i := range conf.Variables {
		varName := strings.TrimSpace(conf.Variables[i].Variable)
		varExpr := strings.TrimSpace(conf.Variables[i].Value)

		if conf.Debug {
			fmt.Printf("==> eval(%q)\n", varExpr)
		}
		val, err := variable.Eval(varExpr, conf.Debug)
		if err != nil {
			return "", err
		}
		if conf.Debug {
			fmt.Println(val)
		}
		buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Importpath, varName, val))
		if i < len(conf.Variables)-1 {
			buf.WriteByte(' ')
		}
	}
	return buf.String(), nil
}
