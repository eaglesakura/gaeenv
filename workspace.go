package gaeenv

import (
	"fmt"
	"os"
	"strings"
)

/*
go testではカレントディレクトリが変動するため、アセットのロードが面倒になる.
それを解決するため、必ず規定ディレクトリになるように調整する.
*/
func updateWorkspacePath(workspacePath string) error {
	if !isDevAppServer() {
		return nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// ".git"ディレクトリを検索する
	for len(pwd) > 0 {
		gitDir := fmt.Sprintf("%v%v.git", pwd, string(os.PathSeparator))
		info, _ := os.Stat(gitDir)
		if info != nil && info.IsDir() {
			// .gitディレクトリを発見
			pwd = fmt.Sprintf("%v%v%v", pwd, string(os.PathSeparator), workspacePath)
			break
		} else {
			// 一個上に戻る
			pwd = pwd[0:strings.LastIndex(pwd, string(os.PathSeparator))]
		}
	}

	// ワークスペースパス再設定
	return os.Chdir(pwd)
}
