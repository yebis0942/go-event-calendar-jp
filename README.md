# golang-jp-event-calendar

connpass上で運営されているGo言語の地域コミュニティのイベント開催情報をまとめます。

## 実行方法

```
# インストール
go install github.com/yebis0942/golang-jp-event-calendar/cmd/golang-jp-event-calendar@latest

# 実行: yyyymmで指定した年月のイベントが{yyymm}.icsにiCalendar形式で出力されます
export CONNPASS_API_KEY=...
golang-jp-event-calendar -yyyymm 200601
```

## 設定

`config.go`に収集対象のconnpassのグループIDを記載してください。
