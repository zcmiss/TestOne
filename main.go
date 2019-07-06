package main

import (
	. "qp/model"
	. "qp/routing"

	"github.com/kataras/iris"
)

func main() {

	app := iris.New()

	// 连接池
	GeneratePool()

	// 路由
	Loading(app)

	//启动服务
	app.Run(iris.Addr(":8899"), iris.WithConfiguration(iris.YAML("./app.yml")))
}
