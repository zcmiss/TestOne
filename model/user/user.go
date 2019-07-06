package user

import (
	. "qp/model/core"
)

// 校验TOKEN
func CheckToken(token string, userId string) bool {
	obj := One(nil, "SELECT token FROM user WHERE id="+userId)
	if _, ok := obj["token"]; ok && (obj["token"] != "") && (obj["token"] == token) {
		return true
	}
	return false
}
