package variable

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/lmika/shellwords"
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
	nameAndArgs := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(expr), "$("), ")")
	if nameAndArgs == "" {
		return "", nil
	}

	output, err := v.exec(nameAndArgs, debug)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (v *CmdVar) exec(nameAndArgs string, debug bool) (output []byte, err error) {
	if runtime.GOOS != "windows" {
		if sh := os.Getenv("SHELL"); sh != "" {
			return v.execByShell(sh, nameAndArgs, debug)
		}
	}
	return v.execByNative(nameAndArgs, debug)
}

func (v *CmdVar) execByNative(nameAndArgs string, debug bool) (output []byte, err error) {
	if debug {
		fmt.Println("==>", nameAndArgs)
	}
	fields := shellwords.Split(nameAndArgs)
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

func (v *CmdVar) execByShell(sh, cmds string, debug bool) (output []byte, err error) {
	if debug {
		fmt.Println("==>", fmt.Sprintf("%s -c %q", sh, cmds))
	}
	return exec.Command(sh, "-c", cmds).Output()
}

// Match 表达式是否可以使用当前变量求值
func (v *CmdVar) Match(expr string) (matched bool) {
	return strings.HasPrefix(expr, "$(") && strings.HasSuffix(expr, ")")
}
