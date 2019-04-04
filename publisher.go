//package publisher
package main

import (
  "context"
  "log"
  "os"

  "cloud.google.com/go/pubsub"
  "github.com/dullgiulio/pingo"
  "github.com/zenkigen/cloud-pubsub-utils/lib"
)

func main() {
  proj := os.Getenv("GOOGLE_PROJECT_ID")
  if proj == "" {
    log.Printf("GOOGLE_PROJECT_ID is not set. ERR:[%v]", os.Stderr)
    os.Exit(1)
  }
  plugin := protocol.Plugin{"HelloPlugin", "./plugins/hello", map[string]string{"name": "Yoshimo"}}
  Do(proj, "test", &plugin)
}

func Do(proj string, topicName string, plugin *protocol.Plugin) {
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

func createMessage(plugin *protocol.Plugin) (string, error) {
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
