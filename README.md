# healthd

## description

スクリプトを定期実行。  
簡易なHTTPServerを起動し、  
指定URLにアクセスすることでスクリプトのstatus codeに応じ、  
結果をJSON形式で返却  

http://[server]:8080/healthcheck

結果

```
{
  "Result": "SUCCESS",
}
```

```
{
  "Result": "FAILED",
}
```

configでメンテナンスファイルを指定しておくとサービス正常稼働時でもFAILEDが返却されます。  
Scriptの実行はintervalで指定した間隔で実行されます。  
REST APIへのアクセスとScriptの実行は非同期になっています。

## Library

```
github.com/ant0ine/go-json-rest/rest
gopkg.in/natefinch/lumberjack.v2
```

## build

```
go build
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
