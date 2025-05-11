# go-event-calendar-jp

connpass上で運営されているGo言語のコミュニティのイベント開催情報をまとめます。

## 設定

`config.go`に収集対象のconnpassのグループIDを記載してください。

## 開発

### `go-event-calendar-jp`コマンド

動作確認用のコマンドです。

`-yyyymm`で指定した年月のイベントが`{yyymm}.ics`にiCalendar形式で出力されます。

```
export CONNPASS_API_KEY=...
go run ./cmd/go-event-calendar-jp -yyyymm 200601
```

### 自動テスト

`make test` を実行してください。

## Inspired by

- [golang.jp ブログ](https://blog.golang.jp/)
- [地域.rbカレンダー](https://sue445.github.io/regional-rb-calendar/)
