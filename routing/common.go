package routing

import (
	. "encoding/json"
	"errors"
	. "fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	. "os"
	. "path/filepath"
	. "qp/config"
	"qp/controller"
	"qp/model/admin"
	. "qp/model/core"
	"qp/model/user"
	. "qp/tool"
	. "reflect"
	"strings"
	. "time"
)

type Object struct {
	handler iris.Handler
}

// 参数校验
func validate(rules iris.Map, params iris.Map) error {

	return nil
}

// 路由映射
func mapping(app *iris.Application, routeName string, method string, needAuth bool, params iris.Map, handler []string) error {
	// 请求方法, 控制器名，操作名
	method, controllerName, actionName := Ufirst(method), handler[0], Ufirst(handler[1])
	// 跨域访问中间件
	corsHandler := Object{
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"}, //允许通过的主机名称
			AllowCredentials: true,
			AllowedHeaders:   []string{"Authorization", "Origin", "Accept", "X-Requested-With", "PLATFORM", "Content-Type"},
			AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			Debug:            false,
		}),
	}
	// 平台校验中间件
	platformHandler := Object{
		func(c iris.Context) {
			// 平台标识
			platform := c.Params().Get("platform")
			//请求时间开始计时
			c.Values().Set("requestCurrentTime", Now().UnixNano()/1e3)
			// 平台标识验证
			if PlatformDbConfigs[platform] != nil {
				CurrentPlatformPool = PlatformPools[platform]
				c.Next()
			} else {
				Output(&c, "没有这个平台号", "无法识别的平台号", ERROR, nil)
			}
		},
	}
	// 用户认证中间件
	authenticateHandler := Object{
		func(c iris.Context) {
			if needAuth {
				token, userId, ok := JwtToken(c)
				if !(ok && ((IsBackendService && admin.CheckToken(token, userId)) || (IsMobileService && user.CheckToken(token, userId)))) {
					return
				}
			}
			c.Next()
		},
	}
	// 参数校验中间件
	paramsHandler := Object{
		func(c iris.Context) {
			result := make(map[string]interface{})
			fields := BuildParams(c)
			isTrue, label := true, ""
			for field, tmp := range params {
				if !isTrue {
					continue
				}
				cfg := tmp.(map[string]interface{})
				// 校验规则
				rules := make([]string, 0)
				if _, ok := cfg["rules"]; ok && (TypeOf(cfg["rules"]).String() == "string") && (cfg["rules"] != "") {
					rules = strings.Split(cfg["rules"].(string), " ")
				}
				// 参数是否必填判断
				if _, ok := fields[field]; (!ok) || (fields[field] == nil) || (fields[field] == "") {
					fields[field] = ""
					if _, ok = cfg["default"]; ok {
						fields[field] = cfg["default"]
					}
				}
				// 参数正则判断
				fields[field], isTrue = CheckField(fields[field], rules)
				if !isTrue {
					isTrue = false
					// 字段名称
					if _, ok := cfg["label"]; ok {
						label = cfg["label"].(string)
					}
				} else {
					result[field] = fields[field]
				}
			}
			if !isTrue {
				Output(&c, label+" 校验失败", "请求失败", ERROR, nil)
			} else {
				c.Values().Set("__", result)
				c.Next()
			}
		},
	}
	// 业务执行方法
	operateHandler := Object{
		func(c iris.Context) {
			ctrl := controller.Object{&c}
			_, err := Todo(ctrl, "Action", []interface{}{controllerName, actionName})
			if err != nil {
				Output(&c, err.Error(), "请求失败", ERROR, nil)
			}
		},
	}
	// 加载完成
	_, err := Todo(
		app,
		method,
		[]interface{}{
			routeName,
			corsHandler.handler,
			platformHandler.handler,
			authenticateHandler.handler,
			paramsHandler.handler,
			operateHandler.handler,
		},
	)
	return err
}

// 路由加载
func handle(app *iris.Application, routingName string, routingConfig interface{}) error {
	if (routingConfig == nil) || (TypeOf(routingConfig).String() != "map[interface {}]interface {}") {
		return errors.New("路由[" + routingName + "]的配置数据错误")
	}
	ncfg := routingConfig.(map[interface{}]interface{})
	// 校验请求方式
	method := strings.ToUpper(CheckStringField(ncfg, "method", ""))
	if (method == "") || (strings.Index("POST,GET,OPTIONS,DELETE,PUT,HEAD", method) == -1) {
		return errors.New("路由[" + routingName + "]的请求方式错误")
	}
	// 是否带上平台标识
	withPlatform := CheckBoolField(ncfg, "withPlatform", false)
	if withPlatform {
		routingName = "/{platform:string}/" + routingName
	}
	// 是否需要认证
	needAuth := CheckBoolField(ncfg, "needAuth", false)
	// 参数校验
	params := MapII2MapSI(CheckMapField(ncfg, "params", nil))
	// 操作校验
	handler := strings.Split(strings.ToLower(CheckStringField(ncfg, "handler", "")), ".")
	if len(handler) != 2 {
		return errors.New("路由[" + routingName + "]的执行方法未定义")
	}
	// 打印路由
	bs, _ := Marshal(params)
	Println("  [" + routingName + "][" + method + "] -> [" + handler[0] + "." + handler[1] + "] -> " + string(bs))
	// 路由映射
	return mapping(app, routingName, method, needAuth, params, handler)
}

// 遍历配置
func (routing Object) ScanRoutingPath(app *iris.Application, routingPath string) {
	err := Walk(routingPath, func(path string, f FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 文件名称
		filename := f.Name()
		// 文件校验
		if (f == nil) || f.IsDir() {
			return nil
		}
		if !strings.HasSuffix(filename, ".yml") {
			return errors.New("路由文件后缀错误")
		}
		// 解析配置
		routings, err := Yaml(routingPath + filename)
		if err != nil {
			return err
		}
		filename = strings.TrimRight(filename, ".yml")
		Println("模块[" + strings.ToUpper(filename) + "] ->")
		for rName, rConfig := range routings {
			err := handle(app, rName, rConfig)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		Println("[error] 路由配置文件解析错误，" + err.Error() + "....")
		Exit(1)
	}
}

// 初始化路由
func Loading(app *iris.Application) {
	// 服务平台
	platform := "移动端"
	// 路由目录
	routingPath := "./routing/mobile/"
	if IsBackendService {
		platform = "管理端"
		routingPath = "./routing/backend/"
	}
	Println("\n>>>>  初始化[" + platform + "]路由配置....")
	// 执行加载路由配置
	_, err := Todo(Object{}, "ScanRoutingPath", []interface{}{app, routingPath})
	if err != nil {
		Println(err.Error())
	}
	Println()
}
