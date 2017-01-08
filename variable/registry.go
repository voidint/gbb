package variable

import "errors"

var (
	// ErrExprEmpty 表达式为空错误
	ErrExprEmpty = errors.New("express is empty")
	// ErrExpr 非法表达式
	ErrExpr = errors.New("invalid express")

	// ErrCmdEmpty 命令为空错误
	ErrCmdEmpty = errors.New("command is empty")
	// ErrCmd 非法命令
	ErrCmd = errors.New("invalid command")
)

// Variabler 变量接口
type Variabler interface {
	// Eval 根据表达式(shell命令、golang模板表达式...)求值
	Eval(expr string, debug bool) (val string, err error)
	// Match 按照表达式检查是否匹配该种类型变量
	Match(expr string) (matched bool)
}

// builtinVars 内建变量集合
var builtinVars = []Variabler{
	DefaultDateVar,
	DefaultGitCommitVar,
	NewCmdVar(),
}

// Eval 逐一用当前内建的变量对表达式求值。
// 若内建变量无一匹配表达式，则返回ErrExpr。
func Eval(expr string, debug bool) (val string, err error) {
	for i := range builtinVars {
		if builtinVars[i].Match(expr) {
			return builtinVars[i].Eval(expr, debug)
		}
	}
	return "", ErrExpr
}
