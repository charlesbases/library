package xerror

// commom ecode
var (
	StatusOK = add(0, "StatusOK")

	RequestErr   = add(-400, "请求错误")
	Unauthorized = add(-401, "未认证")
	TokenExpires = add(-402, "认证过期")
	AccessDenied = add(-403, "权限不足")
	NothingFound = add(-404, "路由错误")
	ParamInvalid = add(-405, "参数错误")

	ServerErr          = add(-500, "服务器错误")
	ServiceUnreachable = add(-501, "服务暂不可用")
	ServiceTimeout     = add(-502, "服务调用超时")
	DatabaseErr        = add(-505, "数据库操作失败")
	DataNotFound       = add(-506, "未查询到数据")
)
