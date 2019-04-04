package pubsubJobExec

import (
  "context"
  "log"
  "os"
  "sync"
  "time"

  "cloud.google.com/go/pubsub"
  "google.golang.org/api/iterator"
  "github.com/dullgiulio/pingo"
  "github.com/zenkigen/cloud-pubsub-utils/lib"
)

type Subscriber struct {
  PluginDir string
}

func (s *Subscriber) Do(proj string, topicName string, subName string, concurrency int) {
  ctx := context.Background()
  client, err := pubsub.NewClient(ctx, proj)

  // create a new topic if not exist
  topic, err := createTopicIfNotExists(client, topicName)
  if err != nil {
    os.Exit(1)
  }

  // create a new subscription if not exists
  sub, err := createSubscriptionIfNotExists(proj, client, subName, topic)
  if err != nil {
    os.Exit(1)
  }
  defer deleteSubscription(client, sub)

  // pull message
  if err := pullMessages(sub, concurrency, s.PluginDir); err != nil {
    os.Exit(1)
  }
}

func createTopicIfNotExists(client *pubsub.Client, topicName string) (*pubsub.Topic, error) {
  ctx := context.Background()
  t := client.Topic(topicName)
  ok, err := t.Exists(ctx)
  if err != nil {
    // TODO: エラー処理追加
    log.Printf("Error[createTopic]: %v", err)
    return nil, err
  }
  if ok {
    log.Printf("Topic exists [%v]", topicName)
    return t, nil
  }
  t, err = client.CreateTopic(ctx, topicName)
  if err != nil {
    // TODO: エラー処理追加
    log.Printf("Error[createTopic]: %v", err)
    return nil, err
  }
  log.Printf("Topic created [%v]", t)
  return t, nil
}

func createSubscriptionIfNotExists(proj string, client *pubsub.Client, subName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
  ctx := context.Background()
  var subs []*pubsub.Subscription
  it := client.Subscriptions(ctx)
  for {
    s, err := it.Next()
    if err == iterator.Done {
      break
    }
    if err != nil {
      log.Printf("Error[createSubscription] %v", err)
      return nil, err
    }
    subs = append(subs, s)
  }
  var sub *pubsub.Subscription
  for _, s := range subs {
    if s.String() == "projects/" + proj + "/subscriptions/" + subName {
      sub = s
    }
  }
  if sub == nil {
    _sub, err := client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
      Topic: topic,
      AckDeadline: 20 * time.Second,
    })
    if err != nil {
      log.Printf("Error[createSubscription] %v", err)
      return nil, err
    }
    sub = _sub
  }
  log.Printf("Subscription created [%v]", sub)
  return sub, nil
}

func pullMessages(sub *pubsub.Subscription, concurrency int, dir string) error {
  ctx := context.Background()
  ch := make(chan []byte, concurrency)
  wg := sync.WaitGroup{}
  defer close(ch)

  cctx, cancel := context.WithCancel(ctx)
  err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
     msg.Ack()
     log.Printf("Got Message: %v", string(msg.Data))
     // StopMessage を受け取ったら Receive を cancel する
     err := protocol.ParseStopMessage(msg.Data)
     if err == nil {
       log.Printf("Received stop command")
       cancel()
       return
     }
     ch <- msg.Data
     log.Printf("Send to goroutine: %v", string(msg.Data))
     wg.Add(1)
     go doPlugin(ch, dir, &wg)
  })
  // エラーに関わらず、goroutine が完了するのを待つ
  wg.Wait()
  if err != nil {
    log.Printf("Error[pullMessage] %v", err)
    log.Printf("Waiting for finishing goroutine...")
    return err
  }
  return nil
}

func deleteSubscription(client *pubsub.Client, sub *pubsub.Subscription) error {
  ctx := context.Background()
  if client == nil || sub == nil {
    log.Printf("No client and subscription")
    return nil
  }
  if err := sub.Delete(ctx); err != nil {
    log.Printf("Error[deleteSubscription] %v", err)
    return err
  }
  log.Printf("Subscription deleted [%v]", sub)
  return nil
}

func doPlugin(ch <-chan []byte, dir string, wg *sync.WaitGroup) {
  v, ok := <-ch
  defer wg.Done()
  if !ok {
    log.Printf("Error[doPlugin] failed to fetch from channel")
    return
  }
  plugin, err := protocol.ParsePluginMessage(v, dir)
  if err != nil {
    log.Printf("Error[doPlugin] unknown message: %v", err)
    return
  }

  p := pingo.NewPlugin("tcp", plugin.Path)
  p.Start()
  defer p.Stop()

  var res string
  err = p.Call(plugin.Name + ".Exec", plugin.Args, &res)
  if err != nil {
    log.Printf("Error[doPlugin] %v", err)
    return
  }
  log.Printf("Plugin done")
  return
}
