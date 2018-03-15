package gaeenv

import (
	"bufio"
	"bytes"
	"io"
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
	buf, err := ioutil.ReadFile(envFileName)
	if err != nil {
		if required {
			// 読み込みに失敗した場合は必須チェック
			return nil, err
		}
		return map[string]string{}, nil
	}

	result := map[string]string{}

	reader := bufio.NewReader(bytes.NewReader(buf))
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return result, nil
			} else {
				return nil, err
			}
		}
		key, value := parseEnvLine(string(line))
		if key != "" {
			result[key] = value
		}
	}
}
