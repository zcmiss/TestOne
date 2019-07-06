package backend

import (
	. "fmt"
	. "github.com/kataras/iris"
	"qp/model/core"
	"qp/tool"
)

type AdminController struct{}

func (ctrl AdminController) List(c *Context) {
	name  := tool.Param(*c, "name").(string)
	admin := core.One(nil, "SELECT * from admin WHERE name="+name)
	Println(admin)
	Println("admin . list")
}
