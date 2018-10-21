package tool

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/lmika/shellwords"
	"github.com/voidint/gbb/config"
	"github.com/voidint/gbb/util"
	"github.com/voidint/gbb/variable"
)

const (
	ldflagsOPT = "-ldflags"
)

// Builder 编译工具
type Builder interface {
	// Build 编译指定目录下的go源码
	Build(dir string) error
}

var (
	// ErrBuildTool 不支持的编译工具错误
	ErrBuildTool = errors.New("unsupported build tool")
)

// Build 根据配置信息，调用合适的编译工具进行编译。
// 若配置的编译工具不在支持的工具范围内，则返回ErrBuildTool错误。
func Build(conf *config.Config, dir string) (err error) {
	defer util.Chdir(dir, conf.Debug) // init work directory

	if err = setupConfig(conf); err != nil {
		return err
	}
	if strings.HasPrefix(conf.Tool, "go ") {
		return NewGoBuilder(*conf).Build(dir)
	} else if strings.HasPrefix(conf.Tool, "gb ") {
		return NewGBBuilder(*conf).Build(dir)
	}
	return ErrBuildTool
}

func setupConfig(conf *config.Config) (err error) {
	if err = setupVars(conf); err != nil {
		return err
	}
	setupTool(conf)
	return nil
}

// setupVars 若定义了变量，则将变量求值后将值重置到Variable的Value属性中。
func setupVars(conf *config.Config) (err error) {
	for i := range conf.Variables {
		if conf.Debug {
			fmt.Printf("==> eval(%q)\n", conf.Variables[i].Value)
		}
		conf.Variables[i].Value, err = variable.Eval(conf.Variables[i].Value, conf.Debug)
		if err != nil {
			return err
		}
		if conf.Debug {
			fmt.Println(conf.Variables[i].Value)
		}
	}
	return nil
}

// setupTool 若定义了变量，则将变量作为ldflags选项的值追加到tool内容中。
// 变量求值应该在调用本函数前完成。
func setupTool(conf *config.Config) {
	var buf bytes.Buffer
	for i := range conf.Variables {
		buf.WriteString(fmt.Sprintf(`-X "%s.%s=%s"`, conf.Importpath, conf.Variables[i].Variable, conf.Variables[i].Value))
		if i < len(conf.Variables)-1 {
			buf.WriteByte(' ')
		}
	}
	ldflags := buf.String()
	if ldflags == "" {
		return
	}

	if !strings.Contains(conf.Tool, ldflagsOPT) {
		conf.Tool = fmt.Sprintf("%s %s '%s'", conf.Tool, ldflagsOPT, ldflags) // 直接增加-ldflags选项及其值
		return
	}

	cmdArgs := shellwords.Split(strings.Replace(conf.Tool, "=", " ", -1))
	for i := range cmdArgs {
		if cmdArgs[i] == ldflagsOPT && i < len(cmdArgs)-1 {
			// 将计算获得的ldflags选项值追加到原-ldflags选项值内容后
			cmdArgs[i+1] = fmt.Sprintf("%s %s", cmdArgs[i+1], ldflags)
			break
		}
	}

	for i := range cmdArgs {
		if strings.Contains(cmdArgs[i], " ") {
			cmdArgs[i] = fmt.Sprintf("'%s'", cmdArgs[i])
		}
	}

	conf.Tool = strings.Join(cmdArgs, " ")
	return
}
