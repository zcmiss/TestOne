package backend

import (
	. "fmt"
	. "github.com/kataras/iris"
	. "qp/config"
	. "qp/model/core"
	. "qp/tool"
)

type AuthController struct{}

func (ac AuthController) Login(c *Context) {
	name := Param(*c, "name").(string)
	admin := Pager(nil, []string{"SELECT * from admin WHERE name='"+name+"'"}, nil)
	Println(admin)
	Println("admin . list")
	Output(c, "", "", SUCCESS, admin)
}
