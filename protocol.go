package pubsubTaskPlugin

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-simplejson"
)

// Plugin struct of information to execute remote plugin
type Plugin struct {
	Name string
	Path string
	Args map[string]string
}

// Message format of pub/sub
// TODO: 構造体を外部に隠蔽する方法を調べる
type Message struct {
	Type    string            `json:"type"`
	Command string            `json:"command"`
	Name    string            `json:"name"`
	BinName string            `json:"binName"`
	Args    map[string]string `json:"args"`
}

// ComposePluginMessage creates message of pub/sub
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

// ParsePluginMessage parses message of pub/sub
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

// CreateStopMessage creates message to stop subscriber
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

// ParseStopMessage parses message of stop from publisher
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
