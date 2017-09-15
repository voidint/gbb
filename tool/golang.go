package tool

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lmika/shellwords"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
)

// GoBuilder go内置编译工具
type GoBuilder struct {
	conf *config.Config
}

// NewGoBuilder 返回go内置编译工具实例
func NewGoBuilder(conf config.Config) *GoBuilder {
	return &GoBuilder{
		conf: &conf,
	}
}

// Build 切换到指定工作目录后调用编译工具编译。
func (b *GoBuilder) Build(rootDir string) (err error) {
	if err = setupConfig(b.conf); err != nil {
		return err
	}

	walkF := util.IsMainPkg
	if b.conf.All {
		walkF = util.IsGoPkg
	}
	paths, err := util.WalkPkgsFunc(rootDir, walkF)
	if err != nil {
		return err
	}
	for i := range paths {
		if err = b.buildDir(paths[i]); err != nil {
			return err
		}
	}
	return nil
}

// 切换到指定工作目录，调用指定的编译工具进行编译。
func (b *GoBuilder) buildDir(dir string) error {
	if err := util.Chdir(dir, b.conf.Debug); err != nil {
		return err
	}

	cmdArgs := shellwords.Split(b.conf.Tool)

	if b.conf.Debug {
		fmt.Print("==> ", cmdArgs[0])
		args := cmdArgs[1:]
		for i := range args {
			if strings.Contains(args[i], " ") {
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
