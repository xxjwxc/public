package errors

import (
	stderrors "errors"
	"fmt"
	"testing"
)

// sentinel 用于验证标准库 errors.Is 能否穿透本包的包装
var sentinel = New("sentinel")

func TestUnwrapStdErrorsIs(t *testing.T) {
	// Wrap -> withStack -> withMessage -> sentinel
	wrapped := Wrap(WithMessage(sentinel, "外层消息"), "再包一层")
	if !stderrors.Is(wrapped, sentinel) {
		t.Error("标准库 errors.Is 未能穿透 Wrap 包装的错误链")
	}

	// WithStack
	if !stderrors.Is(WithStack(sentinel), sentinel) {
		t.Error("标准库 errors.Is 未能穿透 WithStack 包装")
	}

	// WithMessagef
	if !stderrors.Is(WithMessagef(sentinel, "msg %d", 1), sentinel) {
		t.Error("标准库 errors.Is 未能穿透 WithMessagef 包装")
	}

	// 与目标不一致时应返回 false
	other := New("other")
	if stderrors.Is(wrapped, other) {
		t.Error("errors.Is 对不相关的错误应返回 false")
	}
}

func TestUnwrapCause(t *testing.T) {
	wrapped := Wrapf(New("根因"), "包装信息")
	if got := Cause(wrapped).Error(); got != "根因" {
		t.Errorf("Cause 期望『根因』,实际 %q", got)
	}
	if got := wrapped.Error(); got != "包装信息: 根因" {
		t.Errorf("Error() 期望『包装信息: 根因』,实际 %q", got)
	}
}

func test1() error {
	return test2()
}

func test2() error {
	return Wrapf(New("something go wrong"), "自定义消息")
}

func TestErr(t *testing.T) {
	err := test1()
	fmt.Println(fmt.Sprintf("%+v", err))
	err = Cause(err) //获取原始对象
	fmt.Println(fmt.Sprintf("%+v", err))
}

func test11() error {
	return test21()
}

func test21() error {
	return New("something go wrong")
}

func TestErr1(t *testing.T) {
	err := test11()
	fmt.Println(fmt.Sprintf("%+v", err))
	err = Cause(err) //获取原始对象
	fmt.Println(fmt.Sprintf("%+v", err))
}
