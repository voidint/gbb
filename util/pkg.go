package util

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// WalkPkgsFunc 返回指定目录下满足过滤条件的go package路径列表
func WalkPkgsFunc(rootDir string, f FiltePkgFunc) (paths []string, err error) {
	return paths, filepath.Walk(rootDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if !info.IsDir() {
			return nil
		}

		if dname := info.Name(); dname == "vendor" || // skip 'vendor' directory
			(dname == "mod" && strings.HasSuffix(path, filepath.Join("pkg", "mod"))) || // skip '$GOPATH/pkg/mod' directory
			(runtime.GOOS != "windows" && strings.HasPrefix(dname, ".")) { // skip hidden directories
			return filepath.SkipDir
		}

		yes, err := f(path)
		if err != nil {
			return err
		}
		if yes {
			paths = append(paths, path)
		}
		return nil
	})
}

// FiltePkgFunc golang包过滤函数
type FiltePkgFunc func(dir string) (yes bool, err error)

// IsMainPkg 判断指定目录是否是main package
func IsMainPkg(dir string) (yes bool, err error) {
	if dir == "" {
		return false, nil
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return false, err
	}
	_, yes = pkgs["main"]
	return yes, nil
}

// IsGoPkg 判断指定目录是否是go package
func IsGoPkg(dir string) (yes bool, err error) {
	if dir == "" {
		return false, nil
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return false, err
	}
	return len(pkgs) > 0, nil
}
