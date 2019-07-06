package controller

import (
	"github.com/kataras/iris"
	. "qp/config"
	"qp/controller/backend"
	. "qp/tool"
)

type Object struct {
	Ctx *iris.Context
}

var cmapping = map[string]map[string]interface{}{
	"mobile": iris.Map{},
	"backend": iris.Map{
		"index": backend.IndexController{},
		"auth":  backend.AuthController{},
		"admin": backend.AdminController{},
	},
}

func (co Object) Action(controlerName string, actionName string) {
	module := "mobile"
	if IsBackendService {
		module = "backend"
	}
	if _, ok := cmapping[module][controlerName]; ok {
		_, err := Todo(cmapping[module][controlerName], actionName, []interface{}{co.Ctx})
		if err != nil {
			Output(co.Ctx, err.Error(), "请求失败", ERROR, nil)
		}
	} else {
		Output(co.Ctx, "控制器[/controller/"+module+"/"+controlerName+"]不存在", "请求失败", ERROR, nil)
	}
}
