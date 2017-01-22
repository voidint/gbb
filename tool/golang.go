package tool

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/voidint/gbb/config"
)

// GoBuilder go内置编译工具
type GoBuilder struct {
	conf *config.Config
}

// NewGoBuilder 返回go内置编译工具实例
func NewGoBuilder(conf *config.Config) *GoBuilder {
	return &GoBuilder{
		conf: conf,
	}
}

// Build 切换到指定工作目录后调用编译工具编译。
func (b *GoBuilder) Build(rootDir string) error {
	paths, err := walkMainDir(rootDir)
	if err != nil {
		return err
	}
	for i := range paths {
		if err = buildDir(b.conf, paths[i]); err != nil {
			return err
		}
	}
	return nil
}

// 查找root及其子孙目录下所有的main包路径
func walkMainDir(rootDir string) (paths []string, err error) {
	return paths, filepath.Walk(rootDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.IsDir() {
			if info.Name() == "vendor" ||
				(runtime.GOOS != "windows" && strings.HasPrefix(info.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}

		if ext := filepath.Ext(path); ext != ".go" {
			return nil
		}

		yes, err := hasMain(path)
		if err != nil {
			return err
		}
		if yes {
			paths = append(paths, filepath.Dir(path))
		}

		return nil
	})
}

// hasMain 返回指定的go源文件中是否含程序入口。
// Thanks to @stevewang
func hasMain(srcfile string) (bool, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcfile, nil, 0)
	if err != nil {
		return false, err
	}

	if f.Name.Name != "main" {
		return false, nil
	}

	for _, decl := range f.Decls {
		fnDecl, ok := decl.(*ast.FuncDecl)
		if ok && fnDecl.Name != nil && fnDecl.Name.Name == "main" {
			return true, nil
		}
	}
	return false, nil
}

func walkPkgDir(rootDir string) (paths []string, err error) {
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
	path = strings.TrimSpace(path)
	if path == "" {
		return false, nil
	}
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}

	for i := range infos {
		if infos[i].IsDir() {
			continue
		}
		if ext := filepath.Ext(infos[i].Name()); ext == ".go" { // TODO 排除go test目录
			return true, nil
		}
	}
	return false, nil
}
