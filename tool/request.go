package tool

import (
	. "encoding/json"
	. "io/ioutil"
	"strconv"
	"strings"
	"time"

	. "qp/config"

	"github.com/dgrijalva/jwt-go"
	. "github.com/kataras/iris"
)

func BuildParams(c Context) map[string]interface{} {
	fields := make(map[string]interface{})
	request := c.Request()
	// Post表单数据
	request.PostFormValue("")
	for k, v := range request.PostForm {
		fields[k] = v
		if len(v) == 1 {
			fields[k] = v[0]
		}
	}
	// 常规数据
	request.ParseForm()
	for k, v := range request.Form {
		fields[k] = v
		if len(v) == 1 {
			fields[k] = v[0]
		}
	}
	// JSON数据
	bs, _ := ReadAll(request.Body)
	if len(bs) > 0 {
		tmp := make(map[string]interface{})
		err := Unmarshal(bs, &tmp)
		if err == nil {
			for k, v := range tmp {
				fields[k] = v
			}
		}
	}
	// 头部数据
	for k, v := range request.Header {
		fields[k] = v
		if len(v) == 1 {
			fields[k] = v[0]
		}
	}
	return fields
}

func Param(c Context, key string) interface{} {
	tmp := c.Values().Get("__")
	if (tmp == nil) || (!CheckIsMap(tmp)) {
		return ""
	}
	params := tmp.(map[string]interface{})
	if _, ok := params[key]; ok {
		return params[key]
	}
	return ""
}

func JwtToken(c Context) (string, string, bool) {
	tokenStr := strings.Replace(c.GetHeader("Authorization"), "bearer ", "", 7)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenKey), nil
	})
	if (err != nil) || (token == nil) {
		Output(&c, "Token不存在, 或者解板token失败", "需要令牌", TOKEN_EMPTY, nil)
		return "", "", false
	}
	//解码-提交的TOKEN
	claim, _ := token.Claims.(jwt.MapClaims)
	//用户ID
	userId := claim["id"].(string)
	_, err = strconv.Atoi(userId)
	if err != nil {
		Output(&c, "解析Token-ID出错", "令牌格式有误", TOKEN_EMPTY, nil)
		return "", "", false
	}
	//创建时间
	created := int(claim["created"].(float64))
	//当前时间
	current := NowTime()
	//有效时间
	expired := MobileTokenExpire
	if IsBackendService {
		expired = BackendTokenExpire
	}
	//TOKEN已过期
	if current-created > expired {
		Output(&c, "Token与保存的数据不一致", "令牌已经过期", TOKEN_EXPIRED, nil)
		return "", "", false
	}
	return tokenStr, userId, true
}

func Output(c *Context, internalMsg string, clientMsg string, code int16, data interface{}) {
	currentTime := time.Now().UnixNano() / 1e3
	diffTime := (*c).Values().GetInt64Default("requestCurrentTime", currentTime)
	timeConsumed := currentTime - diffTime
	if data == nil {
		data = make(map[string]interface{})
	}
	// clientMsg -> 展示给客户看的
	// internalMsg -> 内部错误给程序员看的,开发和定位问题的时候特别有用
	// timeConsumed -> 总耗时，毫秒
	// data -> 返回数据
	result := Map{"code": code, "clientMsg": clientMsg, "internalMsg": internalMsg, "timeConsumed": timeConsumed, "data": data}
	bs, _ := Marshal(result)
	(*c).Header("Content-Type", "application/json")
	if (len(bs) > 5*1024) && (*c).ClientSupportsGzip() {
		(*c).WriteGzip(bs)
	} else {
		(*c).Write(bs)
	}
}
