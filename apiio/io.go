package apiio

// 输出
func out(code int, data interface{}, msg string) Package {
	if data == nil {
		data = struct{}{}
	}
	return Package{
		code,
		data,
		msg,
	}

}

// 成功
func Success(data interface{}, msg string) Package {
	if msg == "" {
		msg = "ok"
	}
	return out(OK, data, msg)
}

// 失败
func Fail(code int, msg string, data interface{}) Package {
	return out(code, data, msg)
}
