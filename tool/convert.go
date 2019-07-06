package tool

import (
	. "reflect"
	"strings"
)

func Ufirst(str string) string {
	tmp := ""
	vv := []rune(strings.ToLower(str))
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			vv[i] -= 32
		}
		tmp += string(vv[i])
	}
	return tmp
}

func MapII2MapSI(a interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if (a != nil) && (TypeOf(a).String() == "map[interface {}]interface {}") {
		for k, v := range a.(map[interface{}]interface{}) {
			if (TypeOf(k).String() == "string") && (k.(string) != "") {
				if (v != nil) && (TypeOf(v).String() == "map[interface {}]interface {}") {
					v = MapII2MapSI(v)
				}
				result[k.(string)] = v
			}
		}
	}
	return result
}

func MapSS2MapSI(a interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if (a != nil) && (TypeOf(a).String() == "map[string]string") {
		for k, v := range a.(map[string]string) {
			result[k] = v
		}
	}
	return result
}
