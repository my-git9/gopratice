package errs

import (
	"errors"
	"fmt"
)

// 中心式 error 易于演进

var (
	ErrPointerOnly = errors.New("orm: 只支持指向一级指针的结构体")
)

// @ErrUnsupportedExpression 40001 原因是你输入了不支持的类型
// 解决方案：使用正确的类型
func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm: 不支持的表达式类型: %v", expr)
}

func NewErrUnknownField(field any) error {
	return fmt.Errorf("orm: 未知字段: %v", field)
}

func NewErrNoCondition() error {
	return fmt.Errorf("orm: 没有条件")
}

func NewErrUnknownModel() error {
	return fmt.Errorf("orm: 未知模型")
}
