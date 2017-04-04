package tool

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/voidint/gbb/config"
)

// GBBuilder gb编译工具
type GBBuilder struct {
	conf *config.Config
}

// NewGBBuilder 返回gb编译工具实例
func NewGBBuilder(conf *config.Config) *GBBuilder {
	return &GBBuilder{
		conf: conf,
	}
}

// Build 切换到指定工作目录后调用编译工具编译。
func (b *GBBuilder) Build(rootDir string) error {
	return b.buildDir(rootDir)
}

// 切换到指定工作目录，调用指定的编译工具进行编译。
func (b *GBBuilder) buildDir(dir string) error {
	if err := chdir(dir, b.conf.Debug); err != nil {
		return err
	}

	cmdArgs := strings.Fields(b.conf.Tool) // go install ==> []string{"go", "install"}

	flags, err := ldflags(b.conf)
	if err != nil {
		return err
	}
	if flags != "" {
		cmdArgs = append(cmdArgs, "-ldflags", flags)
	}

	if b.conf.Debug {
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
