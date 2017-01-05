package variable

import (
	"time"
)

const (
	// DefaultDateExpr 默认日期变量表达式
	DefaultDateExpr = "{{.Date}}"
)

var (
	// DefaultDateVar 默认的日期变量实例
	DefaultDateVar = NewDateVar(time.RFC3339, DefaultDateExpr)
)

// DateVar 日期变量
type DateVar struct {
	layout string
	expr   string
}

// NewDateVar 实例化日期变量
func NewDateVar(layout, expr string) *DateVar {
	return &DateVar{
		layout: layout,
		expr:   expr,
	}
}

// Eval 表达式变量求值
func (v *DateVar) Eval(_ string, _ bool) (val string, err error) {
	return time.Now().Format(v.layout), nil
}

// Match 表达式是否可以使用当前变量求值
func (v *DateVar) Match(expr string) (matched bool) {
	return v.expr == expr
}
