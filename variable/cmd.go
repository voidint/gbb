package variable

import (
	"fmt"
	"os/exec"
	"strings"
)

// CmdVar 命令变量
type CmdVar struct {
}

// NewCmdVar 实例化命令变量
func NewCmdVar() *CmdVar {
	return &CmdVar{}
}

// Eval 表达式变量求值
func (v *CmdVar) Eval(expr string, debug bool) (val string, err error) {
	nameAndArgs := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(expr), "$("), ")")
	if nameAndArgs == "" {
		return "", nil
	}
	if debug {
		fmt.Println("==>", nameAndArgs)
	}
	output, err := v.exec(nameAndArgs)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (v *CmdVar) exec(nameAndArgs string) (output []byte, err error) {
	fields := strings.Fields(nameAndArgs)
	var cmd *exec.Cmd
	if len(fields) == 1 {
		cmd = exec.Command(fields[0])
	} else if len(fields) > 1 {
		cmd = exec.Command(fields[0], fields[1:]...)
	} else {
		panic("unreachable")
	}

	return cmd.Output()
}

// Match 表达式是否可以使用当前变量求值
func (v *CmdVar) Match(expr string) (matched bool) {
	return strings.HasPrefix(expr, "$(") && strings.HasSuffix(expr, ")")
}
