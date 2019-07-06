package admin

import (
	. "qp/model/core"
)

// 校验TOKEN
func CheckToken(token string, userId string) bool {
	admin := One(nil, "SELECT token FROM admins WHERE id="+userId)
	if _, ok := admin["token"]; ok && (admin["token"] != "") && (admin["token"] == token) {
		return true
	}
	return false
}
