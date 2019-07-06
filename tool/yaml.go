package tool

import (
	"io/ioutil"

	"github.com/kataras/iris"
	"gopkg.in/yaml.v1"
)

func Yaml(filepath string) (iris.Map, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	result := make(iris.Map)
	err = yaml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
