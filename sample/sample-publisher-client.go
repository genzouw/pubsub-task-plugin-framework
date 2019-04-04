package main

import (
  "log"
  "os"

  "github.com/zenkigen/cloud-pubsub-utils"
)

func main() {
  proj := os.Getenv("GOOGLE_PROJECT_ID")
  if proj == "" {
    log.Printf("GOOGLE_PROJECT_ID is not set. ERR:[%v]", os.Stderr)
    os.Exit(1)
  }

  publisher := &pubsubJobExec.Publisher{}
  plugin, err := publisher.NewPlugin("HelloPlugin", "./sample-plugins/hello", map[string]string{"name": "Yoshimo"})
  if err != nil {
    log.Printf("Error[main] %v", err)
    os.Exit(1)
  }
  publisher.Do(proj, "test-topic", plugin)
}
