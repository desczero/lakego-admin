package dsa

import (
    cryptobin_tool "github.com/deatil/go-cryptobin/tool"
)

// 添加错误
func (this DSA) AppendError(err ...error) DSA {
    this.Errors = append(this.Errors, err...)

    return this
}

// 获取错误
func (this DSA) Error() *cryptobin_tool.Errors {
    return cryptobin_tool.NewError(this.Errors...)
}
