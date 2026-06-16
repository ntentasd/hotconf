package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ntentasd/hotconf"
)

type Config struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Debug    bool   `json:"debug"`
	LogLevel string `json:"log_level"`
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	config := newConfig()
	if err = json.NewEncoder(tmpFile).Encode(config); err != nil {
		log.Fatal(err)
	}

	watcher, err := hotconf.NewWatcher(time.Second * 3)
	if err != nil {
		log.Fatal(err)
	}

	if err = watcher.Watch(tmpFile.Name(), handleConfigChange); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = watcher.Start(ctx)
	defer watcher.Stop()

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("running...")
			case <-ctx.Done():
				return
			}
		}
	}()

	<-sigc
}

func newConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     8080,
		Debug:    false,
		LogLevel: "info",
	}
}

func handleConfigChange(path string) {
	c, err := hotconf.Load[Config](path, json.Unmarshal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing config: %v", err)
	}

	fmt.Println("new config detected:")
	fmt.Printf("host: %s\n", c.Host)
	fmt.Printf("port: %d\n", c.Port)
	fmt.Printf("debug: %t\n", c.Debug)
	fmt.Printf("log_level: %s\n", c.LogLevel)
}
