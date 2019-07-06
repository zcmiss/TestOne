package model

import (
	. "fmt"
	. "os"
	. "qp/config"
	. "qp/model/core"
	. "time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
)

//注册数据库连接池
func buildPool(dbAlias string, dbConfig DbConfig) {
	db, err := xorm.NewEngine(dbConfig.DriveName, dbConfig.DriveDsn)
	if err != nil {
		Println("[error] 数据库配置错误," + err.Error() + "....")
		Exit(1)
	}
	// 如果是调试模式，就打印sql语句
	if IsDeveloper {
		db.ShowExecTime(true)
		db.ShowSQL(true)
	}
	err = db.Ping()
	if err != nil {
		Println("[error] [" + dbAlias + "]Mysql数据库连接池创建失败....")
		Exit(1)
	}
	db.SetMaxIdleConns(dbConfig.MaxIdle)
	db.SetMaxOpenConns(dbConfig.MaxConn)
	// 连接生命周期
	db.SetConnMaxLifetime(Duration(9 * Second))
	// 主库单独存放，平台库数组存放
	if dbAlias == "MAIN" {
		MainPool = db
	} else {
		PlatformPools[dbAlias] = db
	}
}

// 创建数据库连接池
func GeneratePool() {
	Println(">>>>  创建数据库连接池....")
	// 平台库
	for dbAlias, dbConfigs := range PlatformDbConfigs {
		platformDbConfig := dbConfigs.(iris.Map)["mysql"].(DbConfig)
		buildPool(dbAlias, platformDbConfig)
	}
	// 主库
	buildPool("MAIN", MainDbConfig)
}
