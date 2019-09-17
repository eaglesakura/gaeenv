package gaeenv

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)

/*
app.yamlを解析し、env_variablesの値を規定値としてセットアップする.

すでにセットアップされている環境変数が優先される.
*/
func parseAppYaml(yamlFileName string) error {
	file, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return err
	}

	var parsedYaml map[string]interface{}
	err = yaml.Unmarshal(file, &parsedYaml)
	if err != nil {
		return err
	}

	envList, ok := parsedYaml["env_variables"]
	if !ok {
		// envListが存在しない
		return nil
	}

	keyValue, ok := envList.(map[interface{}]interface{})
	if !ok {
		// 想定と異なる型だったので落とす
		return fmt.Errorf("cast failed type[%v]\n%v", reflect.TypeOf(envList), keyValue)
	}

	// Key-Valueを取り出す
	for key, value := range keyValue {
		envKey := fmt.Sprintf("%v", key)
		envValue := fmt.Sprintf("%v", value)

		// OSから環境数を取り出す
		// 環境変数が存在したらOSの値を使い、そうでないならapp.yamlを優先する
		osValue, ok := os.LookupEnv(envKey)
		if !ok {
			// 環境変数が存在しないので、app.yaml値を優先
			envCache[envKey] = envValue
			_ = os.Setenv(envKey, envValue)
		} else {
			// 環境変数が存在するので、環境変数を優先
			envCache[envKey] = osValue
		}
	}
	return nil
}
