package pubsubJobExec

import (
  "context"
  "errors"
  "log"
  "os"

  "cloud.google.com/go/pubsub"
  "github.com/dullgiulio/pingo"
)

type Publisher struct {}

type Plugin struct {
  Name string
  Path string
  Args map[string]string
}

func (p *Publisher) NewPlugin(name string, path string, args map[string]string) (*Plugin, error) {
  if name == "" || path == "" {
    return nil, errors.New("empty name or path of plugin")
  }
  return &Plugin{name, path, args}, nil
}

func (p *Publisher) Do(proj string, topicName string, plugin *Plugin) {
  ctx := context.Background()
  client, err := pubsub.NewClient(ctx, proj)

  // create a new topic if not exists
  topic, err := createTopicIfNotExits(client, topicName)
  if err != nil {
    os.Exit(1)
  }

  // create a pub/sub message from plugin
  msg, err := createMessage(plugin)
  if err != nil {
    os.Exit(1)
  }

  // publish a message
  err = publish(client, topic, msg)
  if err != nil {
    os.Exit(1)
  }
}

func createTopicIfNotExits(client *pubsub.Client, topicName string) (*pubsub.Topic, error) {
  ctx := context.Background()
  t := client.Topic(topicName)
  ok, err := t.Exists(ctx)
  if err != nil {
    log.Printf("Error[createTopic]: %v", err)
    return nil, err
  }
  if ok {
    return t, nil
  }
  t, err = client.CreateTopic(ctx, topicName)
  if err != nil {
    log.Printf("Error[createTopic]: %v", err)
    return nil, err
  }
  log.Printf("Topic created[%v]", t)
  return t, nil
}

func publish(client *pubsub.Client, topic *pubsub.Topic, msg string) error {
  ctx := context.Background()
  res := topic.Publish(ctx, &pubsub.Message{Data: []byte(msg)})
  id, err := res.Get(ctx)
  if err != nil {
    log.Printf("Error[publish]: %v", err)
    return err
  }
  log.Printf("Published message; ID: %v", id)
  return nil
}

func createMessage(plugin *Plugin) (string, error) {
  p := pingo.NewPlugin("tcp", plugin.Path)
  p.Start()
  defer p.Stop()

  var res string
  err := p.Call(plugin.Name + ".CreateMessage", plugin.Args, &res)
  if err != nil {
    log.Printf("Error[getMessage] %v", err)
    return "", err
  }
  log.Printf("Message from plugin: %v", res)
  return res, nil
}