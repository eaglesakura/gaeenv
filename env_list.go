package gaeenv

import (
	"bytes"
	"github.com/joho/godotenv"
	"io/ioutil"
	"strings"
)

/*
1行のテキストを解析する
*/
func parseEnvLine(line string) (key, value string) {
	line = strings.Trim(line, " ")
	if len(line) == 0 || line[0:1] == "#" {
		return "", ""
	}

	index := strings.Index(line, "=")
	if index < 0 {
		key = line
		value = ""
	} else {
		key = line[0:index]
		value = line[index+1:]
	}

	return key, value
}

/*
/private/envCache.list が存在する場合、解析して環境変数として登録する.
この関数はパースのみを行い、適用方法は呼び出しに任せる.
*/
func parseEnvList(required bool, envFileName string) (map[string]string, error) {
	result := map[string]string{}
	raw, err := ioutil.ReadFile(envFileName)
	if err == nil {
		result, err = godotenv.Parse(bytes.NewReader(raw))
	}
	if err != nil {
		if required {
			return nil, err
		} else {
			return result, nil
		}
	}
	return result, nil
}
