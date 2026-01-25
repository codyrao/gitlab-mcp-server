package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/codyrao/gitlab-mcp-server/internal/config"
	"github.com/codyrao/gitlab-mcp-server/internal/gitlab"
	"github.com/codyrao/gitlab-mcp-server/internal/mcp"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	transportType := flag.String("transport", "stdio", "Transport type: stdio, sse")
	port := flag.Int("port", 8080, "Port for SSE transport")
	host := flag.String("host", "0.0.0.0", "Host for SSE transport")
	flag.Parse()

	cfg := &config.Config{}
	var err error

	if _, err := os.Stat(*configPath); err == nil {
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			log.Printf("Warning: Failed to load config file: %v", err)
		}
	}

	cfg.Server.Transport = *transportType
	cfg.Server.Host = *host
	cfg.Server.Port = *port

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("Starting GitLab MCP Server...")
	log.Printf("Transport: %s", cfg.Server.Transport)
	log.Printf("GitLab Host: %s", cfg.GitLab.Host)

	client, err := gitlab.NewGitLabClient(cfg.GitLab.Host, cfg.GitLab.Token)
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	server := mcp.NewServer(cfg, client)

	if err := server.Start(ctx, cfg.Server.Transport); err != nil {
		if err != context.Canceled {
			log.Printf("Server error: %v", err)
		}
	}

	log.Println("Server stopped")
}
