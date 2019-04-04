package protocol

import (
  "encoding/json"
)

func ComposePluginMessage(name string, binName string, args map[string]string) (string, error) {
  msg := map[string]interface{}{
    "type": "plugin",
    "name": name,
    "binName": binName,
    "args": args,
  }
  v, err := json.Marshal(msg)
  if err != nil {
    return "", err
  }
  return string(v), nil
}

func CreateStopMessage() (string, error) {
  msg := map[string]interface{}{
    "type": "command",
    "command": "stop",
  }
  v, err := json.Marshal(msg)
  if err != nil {
    return "", err
  }
  return string(v), nil
}
