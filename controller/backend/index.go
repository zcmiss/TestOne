package backend

import (
	"fmt"

	"github.com/kataras/iris"
)

type IndexController struct{}

func (ctrl IndexController) Index(c *iris.Context) {
	fmt.Println("index . index")
}
