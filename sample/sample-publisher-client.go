package samplePublisherClient

import (
	"log"
	"os"

	"github.com/zenkigen/pubsub-task-plugin-framework"
)

func main() {
	proj := os.Getenv("GOOGLE_PROJECT_ID")
	if proj == "" {
		log.Printf("GOOGLE_PROJECT_ID is not set. ERR:[%v]", os.Stderr)
		os.Exit(1)
	}

	publisher := &pubsubTaskPlugin.Publisher{}
	plugin, err := publisher.NewPlugin("HelloPlugin", "./sample-plugins/hello", map[string]string{"name": "Yoshimo"})
	if err != nil {
		log.Printf("Error[main] %v", err)
		os.Exit(1)
	}
	publisher.Do(proj, "test-topic", plugin)
}
