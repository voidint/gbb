package variable

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	// GitCmd 获取最近一次git commit id函数
	GitCmd = `git rev-parse HEAD` //git log --pretty=format:"%H" -1
	// DefaultGitCommitExpr 默认git commit变量表达式
	DefaultGitCommitExpr = "{{.GitCommit}}"
)

var (
	// DefaultGitCommitVar 默认git commit变量
	DefaultGitCommitVar, _ = NewGitCommitVar(GitCmd, DefaultGitCommitExpr)
)

// GitCommitVar git commit变量
type GitCommitVar struct {
	cmd  string
	expr string
}

// NewGitCommitVar 实例化git commit变量
func NewGitCommitVar(cmd, expr string) (v *GitCommitVar, err error) {
	if expr == "" {
		return nil, ErrExprEmpty
	}
	if cmd == "" {
		return nil, ErrCmdEmpty
	}

	return &GitCommitVar{
		cmd:  cmd,
		expr: expr,
	}, nil
}

// Eval 表达式变量求值
func (v *GitCommitVar) Eval(_ string, debug bool) (val string, err error) {
	return v.headCommit(debug)
}

// Match 表达式是否可以使用当前变量求值
func (v *GitCommitVar) Match(expr string) (matched bool) {
	return v.expr == expr
}

func (v *GitCommitVar) headCommit(debug bool) (commit string, err error) {
	var cmd *exec.Cmd
	cmdAndArgs := strings.Fields(v.cmd)
	if len(cmdAndArgs) == 1 {
		cmd = exec.Command(cmdAndArgs[0])
	} else if len(cmdAndArgs) >= 2 {
		cmd = exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...)
	} else {
		panic("unreachable")
	}
	if debug {
		fmt.Println("==>", v.cmd)
	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
