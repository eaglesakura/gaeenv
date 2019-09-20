package gaeenv

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
)

func applyAppengineEnvironments(developMode bool, parsedYaml map[string]interface{}) error {
	if !developMode {
		// load by AppEngine
		envCache["GAE_VERSION"] = os.Getenv("GAE_VERSION")
		envCache["GAE_INSTANCE"] = os.Getenv("GAE_INSTANCE")
		envCache["GAE_SERVICE"] = os.Getenv("GAE_SERVICE")
	} else {
		// insert mock.
		service, ok := parsedYaml["service"]
		if !ok {
			service = "default"
		}

		envCache["GAE_VERSION"] = "__GAE_VERSION__"
		envCache["GAE_INSTANCE"] = "__GAE_INSTANCE__"
		envCache["GAE_SERVICE"] = fmt.Sprintf("%v", service)

		_ = os.Setenv("GAE_VERSION", "__GAE_VERSION__")
		_ = os.Setenv("GAE_INSTANCE", "__GAE_INSTANCE__")
		_ = os.Setenv("GAE_SERVICE", fmt.Sprintf("%v", service))
	}

	fmt.Println("gaeenv.GAE_SERVICE: " + Getenv("GAE_SERVICE"))
	fmt.Println("gaeenv.GAE_VERSION: " + Getenv("GAE_VERSION"))
	fmt.Println("gaeenv.GAE_INSTANCE: " + Getenv("GAE_INSTANCE"))
	return nil
}

func applyEnvVariables(developMode bool, parsedYaml map[string]interface{}) error {
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

/*
app.yamlを解析し、env_variablesの値を規定値としてセットアップする.

すでにセットアップされている環境変数が優先される.
*/
func applyAppYaml(developMode bool, yamlFileName string) error {
	file, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return err
	}

	var parsedYaml map[string]interface{}
	err = yaml.Unmarshal(file, &parsedYaml)
	if err != nil {
		return err
	}

	// apply "env_variables:" block.
	if err := applyEnvVariables(developMode, parsedYaml); err != nil {
		return err
	}

	// apply default environments
	if err := applyAppengineEnvironments(developMode, parsedYaml); err != nil {
		return err
	}

	return nil
}
