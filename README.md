# cloud-pubsub-utils

Cloud pub/sub を利用した並列分散型のタスク実行フレームワークです。
タスクは Plugin で管理することが可能で、publisher から subscriber で実行するプラグインを指定することが出来ます。

# Structure

```
+ publisher.go
+ subscriber.go
+ lib/
   `- protocol.go
+ plugins/
   `- hello.go
```

# How to use

## Prepare

setup environment

```
$ git clone git@github.com:zenkigen/cloud-pubsub-utils.git
$ export GOOGLE_PROJECT_ID="xxxx"
```

get libraries

```
$ go get cloud.google.com/go/pubsub
$ go get google.golang.org/api/iterator
$ go get github.com/dullgiulio/pingo
$ go get github.com/zenkigen/cloud-pubsub-utils
```

# Implement of your plugin

plugin は CreateMessage と Exec の 2 つの関数を実装します。
また、plugin を利用可能にするために、main 関数で plugin 登録を行います。
サンプルとして、publisher で名前を指定して、subscriber で ```hello + ${name}``` を表示する HelloPlugin (./plugins/hello.go) を用視しています。
実装の参考にどうぞ。

## CreateMesaage(args map[string]string, msg *string) error

pub/sub メッセージを生成します。

* args: subscriber で実行する際の引数です (Exec 参照)
* msg: メッセージを格納するポインタです

たいていの場合、サンプルの HelloPlugin (./plugin/hello.go) と同じ実装で問題ないはずです。

## Exec(args map[string]string, res *string) error

subscriber で実行するタスクを実装します。

* args: publisher から指定された引数です
* res: Exec のレスポンスを指定しますが、現状は利用していないので無視して構いません
