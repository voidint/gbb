package tool

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/lmika/shellwords"
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

	if val := Args(shellwords.Split(conf.Tool)).ExtractLdflags(); val != "" {
		buf.WriteString(val)
		buf.WriteByte(' ')
	}

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

// Args 命令行参数
type Args []string

// ExtractLdflags 抽取参数中ldflags所对应的值
func (args Args) ExtractLdflags() string {
	for i, arg := range args {
		if !strings.Contains(arg, "-ldflags") {
			continue
		}
		// eg. go build -ldflags='-w'
		idx := strings.Index(arg, "-ldflags=")
		if idx > -1 {
			return TrimQuotationMarks(arg[idx+len("-ldflags="):])
		}
		if i >= len(args)-1 || !strings.HasSuffix(arg, "-ldflags") {
			return ""
		}
		// eg. go build -ldflags "-w"
		return TrimQuotationMarks(args[i+1])
	}
	return ""
}

// RemoveLdflags 移除ldflags参数及其值
func (args Args) RemoveLdflags() (news Args) {
	for i := range args {
		if !strings.Contains(args[i], "-ldflags") {
			continue
		}
		// eg. go build -ldflags='-w'
		if strings.Contains(args[i], "-ldflags=") {
			args[i] = ""
			continue
		}
		// eg. go build -ldflags "-w"
		if i < len(args)-1 && strings.HasSuffix(args[i], "-ldflags") {
			args[i] = ""
			args[i+1] = ""
			continue
		}
	}

	for i := range args {
		if args[i] == "" {
			continue
		}
		news = append(news, args[i])
	}
	return news
}

// TrimQuotationMarks 去除字符串前后的单/双引号
func TrimQuotationMarks(val string) string {
	if strings.HasSuffix(val, `'`) {
		return strings.TrimPrefix(strings.TrimSuffix(val, `'`), `'`)
	} else if strings.HasSuffix(val, `"`) {
		return strings.TrimPrefix(strings.TrimSuffix(val, `"`), `"`)
	}
	return val
}
