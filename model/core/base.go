package core

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	. "qp/tool"
	. "strconv"
)

// 各平台库连接池
var PlatformPools = make(map[string]*xorm.Engine)

// 当前平台库连接池
var CurrentPlatformPool *xorm.Engine

// 主库连接池
var MainPool *xorm.Engine

// 查询数据
func Search(db *xorm.Engine, sql string, itemCallback func(item map[string]string) (map[string]interface{}, error), otherCallback func(*xorm.Engine, []map[string]interface{}) (interface{}, error)) (interface{}, error) {
	if (db == nil) && (CurrentPlatformPool == nil) {
		return nil, errors.New("数据库不能为空")
	}
	if db == nil {
		db = CurrentPlatformPool
	}
	result := make(map[string]interface{}, 0)
	if sql == "" {
		return result, errors.New("查询语句不能为空")
	}
	list := make([]map[string]interface{}, 0)
	rows, err := db.SQL(sql).QueryString()
	if (err == nil) && (len(rows) > 0) {
		for _, row := range rows {
			if itemCallback != nil {
				item, err := itemCallback(row)
				if err != nil {
					return nil, err
				}
				list = append(list, item)
			} else {
				list = append(list, MapSS2MapSI(row))
			}
		}
	}
	if otherCallback != nil {
		return otherCallback(db, list)
	}
	return list, err
}

// 单个信息
func One(db *xorm.Engine, sql string) map[string]interface{} {
	result, _ := Search(db, sql, nil, func(db *xorm.Engine, list []map[string]interface{}) (interface{}, error) {
		if len(list) > 0 {
			return list[0], nil
		}
		return make(map[string]interface{}), errors.New("信息获取失败")
	})
	return result.(map[string]interface{})
}

// 列表分页
func Pager(db *xorm.Engine, sql []string, itemCallback func(item map[string]string) (map[string]interface{}, error)) map[string]interface{} {
	if len(sql) == 0 {
		return nil
	}
	list := make([]map[string]interface{}, 0)
	tmp, err := Search(db, sql[0], itemCallback, nil)
	if err != nil {
		return nil
	}
	if tmp != nil {
		list = tmp.([]map[string]interface{})
	}
	total := 0
	if len(sql) == 2 {
		tmp, err := Search(db, sql[1], nil, func(db *xorm.Engine, items []map[string]interface{}) (interface{}, error) {
			if len(items) > 0 {
				total := 0
				if _, ok := items[0]["total"]; ok {
					total, _ = Atoi(items[0]["total"].(string))
				}
				return total, nil
			}
			return nil, errors.New("分页信息获取失败")
		})
		if err != nil {
			return nil
		}
		total = tmp.(int)
	}
	result := map[string]interface{}{
		"list":  list,
		"total": total,
	}
	return result
}
