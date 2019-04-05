/**
 * message の protocol を定義するライブラリ
 * protocol の内部実装 (message の形式) を隠蔽する
 * TODO: うまく隠蔽できていないので、デザインパターンを再検討する
 */
package pubsubTaskPlugin

import (
  "errors"
  "encoding/json"

  "github.com/bitly/go-simplejson"
)

// プラグインメッセージを受けた時のプラグイン構造の定義
// プラグインを実行する際に必要な情報を保持
type Plugin struct {
  Name string
  Path string
  Args map[string]string
}

// protocol 内部のメッセージ形式の定義
// TODO: 構造体を外部に隠蔽する方法を調べる
type Message struct {
  Type    string            `json:"type"`
  Command string            `json:"command"`
  Name    string            `json:"name"`
  BinName string            `json:"binName"`
  Args    map[string]string `json:"args"`
}

// [For publisher/plugin]
// subscriber にプラグインを実行させるメッセージを作成する
func ComposePluginMessage(name string, binName string, args map[string]string) (string, error) {
  json := simplejson.New()
  json.Set("type", "plugin")
  json.Set("name", name)
  json.Set("binName", binName)
  json.Set("args", args)
  v, err := json.MarshalJSON()
  if err != nil {
    return "", err
  }
  return string(v), nil
}

// [For subscriber]
// プラグインを実行するメッセージをパースする
func ParsePluginMessage(data []byte, pluginDir string) (*Plugin, error) {
  var msg Message
  err := json.Unmarshal(data, &msg)
  if err != nil {
    return nil, err
  }
  if msg.Type != "plugin" {
    return nil, errors.New("invalid message")
  }
  if len(msg.Name) == 0 || len(msg.BinName) == 0 {
    return nil, errors.New("invalid parameters of plugin message")
  }
  p := Plugin{msg.Name, pluginDir + "/" + msg.BinName, msg.Args}
  return &p, nil
}

// [For publisher]
// subscriber の stop させるコマンドのメッセージを作成する
func CreateStopMessage() (string, error) {
  json := simplejson.New()
  json.Set("type", "command")
  json.Set("command", "stop")
  v, err := json.MarshalJSON()
  if err != nil {
    return "", err
  }
  return string(v), nil
}

// [For subscriber]
// メッセージの受信を stop するコマンドのメッセージをパースする
func ParseStopMessage(data []byte) error {
  var msg Message
  err := json.Unmarshal(data, &msg)
  if err != nil {
    return err
  }
  if msg.Type != "command" {
    return errors.New("invalid message")
  }
  if msg.Command != "stop" {
    return errors.New("invalid parameter of command message")
  }
  return nil
}
