package tool

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func CheckField(field interface{}, rules []string) (interface{}, bool) {
	isTrue := true
	for _, rule := range rules {
		switch rule {
		case "must":
			isTrue = (field != nil) && (field != "")
		case "string":
			isTrue = CheckIsString(field)
		case "int":
			if CheckIsString(field) {
				tmp, err := strconv.Atoi(field.(string))
				if err != nil {
					isTrue = false
				}
				field = tmp
			} else {
				isTrue = CheckIsInt(field)
			}
		case "bool":
			if CheckIsString(field) {
				field = strings.ToUpper(field.(string))
				if (field == "YES") || (field == "TRUE") {
					field = true
				} else if (field == "NO") || (field == "FALSE") {
					field = false
				} else {
					isTrue = false
				}
			} else {
				isTrue = CheckIsBool(field)
			}
		case "float":
			if CheckIsString(field) {
				tmp, err := strconv.ParseFloat(field.(string), 64)
				if err != nil {
					isTrue = false
				}
				field = tmp
			} else {
				isTrue = CheckIsFloat(field)
			}
		case "array":
			isTrue = CheckIsArray(field)
		default:
			if CheckIsString(field) {
				isTrue, _ = regexp.MatchString(``, field.(string))
			} else {
				isTrue = false
			}
		}
		if isTrue == false {
			return nil, false
		}
	}
	return field, true
}

func CheckIsString(a interface{}) bool {
	return strings.HasPrefix(reflect.TypeOf(a).String(), "string")
}

func CheckIsInt(a interface{}) bool {
	return strings.HasPrefix(reflect.TypeOf(a).String(), "int")
}

func CheckIsFloat(a interface{}) bool {
	return strings.HasPrefix(reflect.TypeOf(a).String(), "float")
}

func CheckIsBool(a interface{}) bool {
	return strings.HasPrefix(reflect.TypeOf(a).String(), "bool")
}

func CheckIsMap(a interface{}) bool {
	return strings.HasPrefix(reflect.TypeOf(a).String(), "map[")
}

func CheckIsArray(a interface{}) bool {
	return strings.HasPrefix(reflect.TypeOf(a).String(), "[]")
}

func CheckBoolField(a interface{}, field string, defaultValue bool) bool {
	if (a == nil) || (field == "") {
		return defaultValue
	}
	if reflect.TypeOf(a).String() == "map[string]interface {}" {
		obj := a.(map[string]interface{})
		if _, ok := obj[field]; ok && (reflect.TypeOf(obj[field]).String() == "bool") {
			return obj[field].(bool)
		}
	} else if reflect.TypeOf(a).String() == "map[interface {}]interface {}" {
		obj := a.(map[interface{}]interface{})
		if _, ok := obj[field]; ok && (reflect.TypeOf(obj[field]).String() == "bool") {
			return obj[field].(bool)
		}
	}
	return defaultValue
}

func CheckStringField(a interface{}, field string, defaultValue string) string {
	if (a == nil) || (field == "") {
		return defaultValue
	}
	if reflect.TypeOf(a).String() == "map[string]interface {}" {
		obj := a.(map[string]interface{})
		if _, ok := obj[field]; ok && (reflect.TypeOf(obj[field]).String() == "string") && (obj[field] != "") {
			return obj[field].(string)
		}
	} else if reflect.TypeOf(a).String() == "map[interface {}]interface {}" {
		obj := a.(map[interface{}]interface{})
		if _, ok := obj[field]; ok && (reflect.TypeOf(obj[field]).String() == "string") && (obj[field] != "") {
			return obj[field].(string)
		}
	}
	return defaultValue
}

func CheckMapField(a interface{}, field string, defaultValue map[interface{}]interface{}) map[interface{}]interface{} {
	if (a == nil) || (field == "") {
		return defaultValue
	}
	if reflect.TypeOf(a).String() == "map[string]interface {}" {
		obj := a.(map[string]interface{})
		if _, ok := obj[field]; ok && (reflect.TypeOf(obj[field]).String() == "map[interface{}]interface{}") && (obj[field] != nil) {
			return obj[field].(map[interface{}]interface{})
		}
	} else if reflect.TypeOf(a).String() == "map[interface {}]interface {}" {
		obj := a.(map[interface{}]interface{})
		if _, ok := obj[field]; ok && (reflect.TypeOf(obj[field]).String() == "map[interface {}]interface {}") && (obj[field] != nil) {
			return obj[field].(map[interface{}]interface{})
		}
	}
	return defaultValue
}
