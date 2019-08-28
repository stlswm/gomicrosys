package io

// 网络请求数据包
type Package struct {
	Code int
	Data interface{}
	Msg  string
}
