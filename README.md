# golang-jp-event-calendar

connpass上で運営されているGo言語の地域コミュニティのイベント開催情報をまとめます。

## 設定

`config.go`に収集対象のconnpassのグループIDを記載してください。

## 開発

### `golang-jp-event-calendar`コマンド

動作確認用のコマンドです。

`-yyyymm`で指定した年月のイベントが`{yyymm}.ics`にiCalendar形式で出力されます。

```
export CONNPASS_API_KEY=...
go run ./cmd/golang-jp-event-calendar -yyyymm 200601
```

### 自動テスト

`make test` を実行してください。
