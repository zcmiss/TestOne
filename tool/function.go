package tool

import (
	"errors"
	. "reflect"
)

func Todo(obj interface{}, method string, params []interface{}) (interface{}, error) {
	v := ValueOf(obj)
	m := v.MethodByName(method)
	if m.IsValid() {
		values := make([]Value, 0)
		if (params != nil) && (len(params) > 0) {
			for _, v := range params {
				values = append(values, ValueOf(v))
			}
		}
		return m.Call(values), nil
	}
	return nil, errors.New("[error] 方法[" + method + "]不存在....")
}
