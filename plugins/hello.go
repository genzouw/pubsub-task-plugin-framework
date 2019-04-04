package main

import (
  "errors"

  "github.com/dullgiulio/pingo"
  "github.com/zenkigen/cloud-pubsub-utils/lib"
)

type HelloPlugin struct {}

func (p *HelloPlugin) CreateMessage(args map[string]string, msg *string) error {
  err := checkArgument(args)
  if err != nil {
    return err
  }
  *msg, err = protocol.ComposePluginMessage("HelloPlugin", "hello", args)
  if err != nil {
    return err
  }
  return nil
}

func checkArgument(args map[string]string) error {
  if args == nil || len(args) < 1 {
    return errors.New("no parameters")
  }
  _, ok := args["name"]
  if !ok {
    return errors.New("invalid parameters")
  }
  return nil
}

func main() {
  plugin := &HelloPlugin{}
  pingo.Register(plugin)
  pingo.Run()
}
