package webserver

// commom ecode
var (
	StatusOK = add(0, "StatusOK")

	StatusBadRequest   = add(4000, "请求错误")
	StatusParamInvalid = add(4001, "参数错误")
	StatusUnauthorized = add(4010, "未认证")
	StatusTokenExpires = add(4011, "认证过期")
	StatusAccessDenied = add(4012, "权限不足")
	StatusNotFound     = add(4040, "路由错误")

	StatusServerError        = add(5000, "服务器错误")
	StatusDatabaseErr        = add(5200, "数据库操作失败")
	StatusDataNotFound       = add(5201, "未查询到数据")
	StatusServiceUnreachable = add(5300, "服务暂不可用")
	StatusServiceTimeout     = add(5301, "服务调用超时")
)
