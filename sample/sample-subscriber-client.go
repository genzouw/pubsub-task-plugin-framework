package main

import (
  "log"
  "os"
  "strconv"

  "github.com/zenkigen/cloud-pubsub-utils"
)

func main() {
  proj := os.Getenv("GOOGLE_PROJECT_ID")
  if proj == "" {
    log.Printf("GOOGLE_PROJECT_ID is not set. ERR:[%v]", os.Stderr)
    os.Exit(1)
  }
  concurrency, err := strconv.Atoi(os.Getenv("SUBSCIBER_CONCURRENCY"))
  if err != nil || concurrency == 0 {
    concurrency = 3
  }
  subscriber := &pubsubJobExec.Subscriber{"./sample-plugins"}
  subscriber.Do(proj, "test-topic", "test-subscription", concurrency)
}