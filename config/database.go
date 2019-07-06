package config

import (
	"github.com/kataras/iris"
	. "os"
)

// 数据配置构造
type DbConfig struct {
	AliasName string //数据库别名
	DriveName string //驱动名称
	DriveDsn  string //数据库连接地址
	MaxIdle   int    //最大空闲数
	MaxConn   int    //最大连接数
}

// 主库配置
var MainDbConfig = DbConfig{
	AliasName: "maindb",
	DriveName: "mysql",
	DriveDsn:  "allSysPlatDb:m$re+ghKpmn*sf@tcp(127.0.0.1:3306)/maindb?charset=utf8&multiStatements=true",
	MaxIdle:   5,
	MaxConn:   30,
}

// 平台库配置
var PlatformDbConfigs = iris.Map{
	"CKYX": iris.Map{
		"mysql": DbConfig{
			AliasName: "CKYXMYSQL",
			DriveName: "mysql",
			DriveDsn:  "ckgame:+ghKCkg0&pL*=!@tcp(127.0.0.1:3306)/ckgame?charset=utf8&multiStatements=true",
			MaxIdle:   10,
			MaxConn:   30,
		},
	},
}

// 初始化
func init() {
	if Getenv("IRIS_MODEL") == "release" {
		IsDeveloper = false
		if _, ok := PlatformDbConfigs["QPTEST"]; ok {
			delete(PlatformDbConfigs, "QPTEST")
		}
		CurrentPlatform := Getenv("CURRENT_PLATFORM")
		if _, ok := PlatformDbConfigs[CurrentPlatform]; ok {
			PlatformDbConfigs = iris.Map{
				CurrentPlatform: PlatformDbConfigs[CurrentPlatform],
			}
		}
	}
}
