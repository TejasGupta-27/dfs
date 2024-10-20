package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TejasGupta-27/dfs/internal/node"
	"github.com/TejasGupta-27/dfs/config"
)

func main() {
	cfg := config.LoadConfig()

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