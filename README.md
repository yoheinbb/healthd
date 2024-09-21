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

## Library

```
github.com/ant0ine/go-json-rest/rest
gopkg.in/natefinch/lumberjack.v2
```

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
  "script": "scripts/sample_script",
  "maintenance_file": "/tmp/maintenance",
  "interval": "10s",
  "timeout": "10s"
}
```
