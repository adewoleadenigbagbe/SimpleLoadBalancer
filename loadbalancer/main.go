package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/lb"
)

const (
	configFile = "config.json"
)

func main() {
	lbConfig, err := CreateLBConfig()
	if err != nil {
		log.Fatal(err)
	}

	loadbalancer, err := lb.CreateLB(*lbConfig)
	if err != nil {
		log.Fatal(err)
	}

	//create echo server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", lbConfig.Port),
		Handler: http.HandlerFunc(loadbalancer.Serve),
	}

	// Start server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func CreateLBConfig() (*lb.LbConfig, error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filepath := filepath.Join(currentWorkingDirectory, configFile)
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config lb.LbConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
