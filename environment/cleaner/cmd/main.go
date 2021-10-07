package main

import (
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/environment/cleaner"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal().Err(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err)
	}
	c := cleaner.NewCleaner(clientset, &cleaner.Config{PollInterval: 30 * time.Second})
	if err := c.Run(); err != nil {
		log.Fatal().Err(err)
	}
}
