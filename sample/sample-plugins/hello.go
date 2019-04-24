package main

import (
	"errors"
	"log"

	"github.com/dullgiulio/pingo"
	"github.com/zenkigen/pubsub-task-plugin-framework"
)

// HelloPlugin for plugin
type HelloPlugin struct{}

// CreateMessage creates pub/sub message to execute remote plugin
func (p *HelloPlugin) CreateMessage(args map[string]string, msg *string) error {
	err := checkArgument(args)
	if err != nil {
		return err
	}
	*msg, err = pubsubTaskPlugin.ComposePluginMessage("HelloPlugin", "hello", args)
	if err != nil {
		return err
	}
	return nil
}

// Exec HelloPlugin
func (p *HelloPlugin) Exec(args map[string]string, res *string) error {
	err := checkArgument(args)
	if err != nil {
		return err
	}
	log.Printf("Hello " + args["name"])
	*res = "Hello, " + args["name"]
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
	err := pingo.Run()
	if err != nil {
		log.Fatal(err)
	}
}
