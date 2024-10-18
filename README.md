# healthd

## description

スクリプトを定期実行。  
簡易なHTTPServerを起動し、  
指定URLにアクセスすることでスクリプトの終了ステータスに応じ、  
結果をJSON形式で返却  

port, URLパス、返却文字列はconf/global.jsonに指定
http://[server]:[port]/healthcheck

実行するscriptはconf/script.jsonに指定する
Scriptの実行はintervalで指定した間隔で実行  

結果
exit code 0 の場合

```
{
  "Result": "SUCCESS",
}
```

exit code 1 の場合
```
{
  "Result": "FAILED",
}
```

configでメンテナンスファイルを指定しておくとサービス正常稼働時でもFAILEDが返却される。  
REST APIへのアクセスとScriptの実行は非同期で実行

## build

```
make build
```

## create docker image

```
make image
```

## execute e2e test

ローカルにDockerコンテナを起動  
8080, 8081を利用しsuccess, fail時のe2eテストを行う

```
make e2e
```

## usage

```
./healthd [-global-config=conf/global.json -script-config=conf/script.json]
```
引数指定無しの場合は、下記configがデフォルトで読み込まれます。


```
conf/global.json
conf/script.json
```


## config/global.json

```
{
  "port": "80",
  "urlpath": "healthcheck",
  "ret_success": "SUCCESS",
  "ret_failed": "FAILED"
}
```

## config/script.json

```
{
  "id": "sample",
  "name": "sample",
  "script": "configs/scripts/sample_script",
  "maintenance_file": "/tmp/maintenance",
  "interval": "10s",
  "timeout": "10s"
}
```

## test in local environment with sample config

```
make run-local
```
```
curl localhost/healthcheck
```

## directory structure

```
.
├── Dockerfile
├── Makefile
├── README.md
├── configs   # sample configs
├── e2e       # e2e test
├── go.mod
├── go.sum
├── internal  # internal pkg
└── main.go
```

## refs

### directory structure
* https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages
* https://github.com/golang-standards/project-layout

### code structure
* https://amasuda.xyz/post/2023-01-12-pros-cons-ddd-and-golang/

