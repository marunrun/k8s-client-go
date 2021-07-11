package common

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
)


func ParseYamlToJson(filename string) []byte {
	// 读取YAML
	deployYaml, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// YAML转JSON
	deployJson, err := yaml.ToJSON(deployYaml)
	return deployJson
}