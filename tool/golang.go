package tool

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lmika/shellwords"
	"github.com/voidint/gbb/config"
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

	paths, err := walkPkgs(rootDir)
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
	if err := chdir(dir, b.conf.Debug); err != nil {
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

func walkPkgs(rootDir string) (paths []string, err error) {
	return paths, filepath.Walk(rootDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if !info.IsDir() {
			return nil
		}

		if info.Name() == "vendor" ||
			(runtime.GOOS != "windows" && strings.HasPrefix(info.Name(), ".")) {
			return filepath.SkipDir
		}

		yes, err := isGoPkg(path)
		if err != nil {
			return err
		}
		if yes {
			paths = append(paths, path)
		}
		return nil
	})
}

// isGoPkg 判断路径是否是golang的包
func isGoPkg(path string) (yes bool, err error) {
	if path == "" {
		return false, nil
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return false, err
	}
	return len(pkgs) > 0, nil
}
