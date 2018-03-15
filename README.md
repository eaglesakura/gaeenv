# en) gaeenv Overview

The `gaeenv` library is environment variables utils for GAE/Go.

GAE/Go environment is not equals localserver, UnitTest and GCP Runtime.
If "gaeenv" library use, All environment variables are equals.

 * [All] Environment variables wrote in app.yaml, Those are set to os.Setenv.
 * [All] Environment variables wrote in env.list file(Docker compatible), Those are set to os.Setenv.
 	* Priority
	  1. env.list environments
		2. app.yaml environments
		3. os.Getenv environments
 * [for UnitTest] Current directory change to include app.yaml directory.
 	* Current directory is equals "goapp serve".

# ja) gaeenv Overview

`gaeenv` はGAE/GoでのUnitTest/DevServer/Runtimeそれぞれの環境変数設定・カレントディレクトリ設定の違いを吸収することを目的としている.

GAE/Goの環境変数はローカルサーバー・UnitTest・GCPでの実行でそれぞれ異なる状態で実行されるため、環境変数から読み出しを行うプログラムが動作させづらい場合がある。
gaeenvを使用することで、すべての環境での環境変数やカレントディレクトリ設定を統一し、UnitTestやローカルサーバーでの動作確認をしやすくしたり、開発者個人ごとの設定を反映させやすくなる。

 * [All] app.yamlに記述された環境変数をos.Setenvに渡す.
 * [All] env.list(Docker互換)に記述された環境変数をos.Setenvに渡す.
 	* 優先順
	  1. env.list
		2. app.yaml
		3. os.Getenv
 * [for UnitTest] カレントディレクトリをapp.yamlを配置したディレクトリに変更する.
 	* gaeenv.Init()を呼び出した時点で、 `goapp serve` を行った場合と同じカレントディレクトリに変更される.

# example

```
func init() {
  // call "gaeenv.Init()" function at first.
  gaeenv.Init(&gaeenv.AppengineEnvironment{
    AppYamlPath: "path/to/app.yaml",
    EnvListFiles: []gaeenv.EnvironmentListFile{
      {
        DevOnly:  true,
        Required: false,
        Path:     "../private/env.list",
      },
    },
    RequiredEnvironments: []string{
      "PRIVATE_API_KEY",
    },
  })
}
```
