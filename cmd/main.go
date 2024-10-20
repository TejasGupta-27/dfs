package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dfs/config"
	"dfs/internal/node"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	node, err := node.NewNode(cfg)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	if err := node.Start(); err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	fmt.Println("Node started. Press Ctrl+C to exit.")

	// Wait for interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Shutting down...")
	node.Stop()
}