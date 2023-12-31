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
	pool "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/serverpool"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//create echo server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", lbConfig.Port),
		Handler: http.HandlerFunc(loadbalancer.Serve),
	}

	go loadbalancer.HealthCheck(ctx)

	// Start server
	go func() {
		fmt.Println("server started on port:", lbConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func CreateLBConfig() (*pool.LbConfig, error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filepath := filepath.Join(currentWorkingDirectory, configFile)
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config pool.LbConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
