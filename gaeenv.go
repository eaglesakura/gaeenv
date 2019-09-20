package gaeenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var envCache = map[string]string{}
var errorCache []error

/*
en)
This struct is repository settings holder for GAE/Go.

GAE/Go environment is not equals localserver, UnitTest and GCP Runtime.
If "gaeenv" library use, All environment variables are equals.

 * [All] Environment variables wrote in app.yaml, Those are set to os.Setenv.
 * [All] Environment variables wrote in env.list file(Docker compatible), Those are set to os.Setenv.
 	* Priority
	  1. env.list environments
		2. app.yaml environments
		3. os.Getenv environments
 * [for UnitTest] Current directory change to included app.yaml directory.
 	* Current directory is equals "goapp serve".

ja)
GAE/Goのリポジトリ構成を指定する.
AppEngineの実環境,UnitTest,LocalServer環境の違いを吸収する.
*/
type AppengineEnvironment struct {
	/*
		en)
		This member relative path to app.yaml from repository root directory.
		The repository root path is auto find with ".git" directory.

		ja)
		リポジトリrootから見た、app.yamlへの相対パス.
		リポジトリrootは ".git"ディレクトリの有無によって判定する.

		ex)
		"gae/app.yaml"
	*/
	AppYamlPath string

	/*
		en)
		This runtime is debug(develop, or Unit Test) mode.

		ja)
		実行環境がDebug(開発環境、もしくはUnitTest）であればtrueを指定する
	*/
	DevelopMode bool

	/*
		en)
		This member is relative path to env.list from app.yaml included directory.

		ja)
		dockerの環境変数設定ファイルに沿ったenv.listファイルリスト.
		環境変数を上書きする.
		app.yamlのディレクトリから見た相対パスで記述する点に注意すること.
	*/
	EnvListFiles []EnvironmentListFile

	/*
		en)
		Environment variables names required for GAE/Go server program.
		If environment variable not found or empty, gaeenv.Init() function does return error.

		ja)
		必須で設定されなければならない環境変数一覧.
		もし環境変数が見つからないか空の場合、gaeenv.Init()はエラーコードを返却します。
	*/
	RequiredEnvironments []string
}

/*
en)
This struct is env.list file settings.
env.list file format is Docker compatible.

Repository owner add line "/env.list" into .gitignore, it's better.
Developers private settings(APIKey, Token...) is not equals. If write private settings in "/env.list", Developers can use private settings on development time.

ja)
Dockerで使用するenv.list形式ファイルを環境変数としてアタッチする.
これは開発者ごとの個人環境をセットアップするために使用する.
*/
type EnvironmentListFile struct {
	/*
		en)
		If use env.list file on development time, It should set "true" to this member.

		ja)
		開発環境のみ有効な場合はtrue
	*/
	DevOnly bool

	/*
		en)
		If this member is set to "true" and not found env.list file, "gaeenv.Init()" function does return error.
		However, if "DevOnly && IsDevAppServer()" are "true", It setting is ignored.

		ja)
		読み込み必須の場合はtrue.
		ただし、DevOnly && IsDevAppServer()を満たす場合は無視される.
	*/
	Required bool

	/*
		en)
		This member is relative path to "env.list" file from app.yaml include directory.

		ja)
		Docker互換のenv.listファイルへのパス.
		複数のKeyが記述されている場合、後に記述されたものが優先される.

		ex)
		"../env.list"
	*/
	Path string
}

/*
en)
"GetError" is return cached error at "gaeenv.Init()".
If "gaeenv.Init()" returned nil, this function also returns nil.

ja)
内部でキャッシュされたエラーを返却する.
*/
func GetError() error {
	err := errorCache
	if len(err) > 0 {
		return err[len(err)-1]
	}
	return nil
}

func pushError(err error) error {
	if err != nil {
		errorCache = append(errorCache, err)
	}
	return err
}

/*
環境変数を上書きする.
ただし、app.yamlからロードされた環境変数リストにKeyがなければ panic() として扱う.
*/
func applyEnv(envList map[string]string) error {
	for key, value := range envList {
		_, ok := envCache[key]
		if !ok {
			return fmt.Errorf("error, environment variable key[%v]. not found in app.yaml", key)
		} else {
			envCache[key] = value
			_ = os.Setenv(key, value)
		}
	}
	return nil
}

/*
en)
"GetEnvList" return all environment variables.
Those are initialized values at "gaeenv.Init()".

ja)
構築された環境変数一覧を取得する
*/
func GetEnvList() map[string]string {
	result := map[string]string{}
	for key, value := range envCache {
		result[key] = value
	}
	return result
}

/*
en)
"Getenv" does return environment variable.
If environment key not found, this function call "panic()".
Developer should "recover()" or rebuild programs.

ja)
環境変数をキャッシュから取得する.
*/
func Getenv(key string) string {
	value, ok := envCache[key]
	if !ok {
		panic(fmt.Errorf("error, environment variable key[%v]. not found in app.yaml", key))
	}
	return value
}

/*
en)
"Getenv" does return environment variable. It convert to integer.
If environment key not found, this function call "panic()".
Developer should "recover()" or rebuild programs.

ja)
環境変数をキャッシュから取得する.
*/
func GetenvI(key string) int64 {
	result, err := strconv.ParseInt(Getenv(key), 10, 64)
	if err != nil {
		panic(err)
	}
	return result
}

/*
en)
"Getenv" does return environment variable. It convert to float.
If environment key not found, this function call "panic()".
Developer should "recover()" or rebuild programs.

ja)
環境変数をキャッシュから取得する.
*/
func GetenvF(key string) float64 {
	result, err := strconv.ParseFloat(Getenv(key), 64)
	if err != nil {
		panic(err)
	}
	return result
}

/*
en)
"Getenv" does return environment variable. It convert to float.
If environment key not found, this function does return empty string and "ok=false".

ja)
環境変数を取得する.
環境変数が見つからない場合、この関数は空文字とfalseを返却する.
*/
func Lookenv(key string) (value string, ok bool) {
	value, ok = envCache[key]
	if !ok {
		return "", ok
	}
	return
}

/*
en)
This function does initialize environment variables and current directory.
If this function failed to call, does return error.
Developers may error handle on async. It should use "gaeenv.GetError()" function.

ja)
実行環境をセットアップする.
セットアップに失敗した場合errorを返却するが、非同期的に取得するために "gaeenv.GetError()" でも同じエラーを返却する.
*/
func Init(repo *AppengineEnvironment) error {
	yamlDirectory := ""
	var yamlFileName string
	if strings.LastIndex(repo.AppYamlPath, "/") > 0 {
		yamlFileName = repo.AppYamlPath[(strings.LastIndex(repo.AppYamlPath, "/") + 1):]
		yamlDirectory = repo.AppYamlPath[:strings.LastIndex(repo.AppYamlPath, "/")]
	} else {
		yamlFileName = repo.AppYamlPath
	}

	// カレントディレクトリを指定
	if repo.DevelopMode {
		if err := pushError(updateWorkspacePath(yamlDirectory)); err != nil {
			return err
		}
	}

	// app.yamlからデフォルトの環境変数を設定
	if err := pushError(applyAppYaml(repo.DevelopMode, yamlFileName)); err != nil {
		return err
	}

	// env.listのセットアップ
	for _, file := range repo.EnvListFiles {
		if file.DevOnly && !repo.DevelopMode {
			// 開発環境のみセットアップする
			continue
		}

		// envCache.listから、環境ごとのセットアップ
		list, err := parseEnvList(file.Required, file.Path)
		if pushError(err) != nil {
			return err
		}
		if err = pushError(applyEnv(list)); err != nil {
			return err
		}
	}

	// 必須環境変数をチェック
	for _, key := range repo.RequiredEnvironments {
		value, ok := Lookenv(key)
		if !ok || value == "" {
			return pushError(fmt.Errorf("environment key[%v] is empty. edit 'app.yaml' or 'env.list' or Environment Variables", key))
		}
	}

	return nil
}
