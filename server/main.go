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
	"time"

	database "github.com/adewoleadenigbagbe/simpleloadbalancer/server/db"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/server/handlers"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
)

var (
	ServerConfigFile = "config.json"
)

type serverConfig struct {
	Ip       string
	Port     int
	Protocol string
}

func CreateServerConfig() serverConfig {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filepath := filepath.Join(currentWorkingDirectory, ServerConfigFile)
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	var config serverConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {
	var err error

	//connect to db
	database.ConnectToSqlite()

	//new config
	config := CreateServerConfig()

	e := echo.New()
	e.Logger.SetLevel(echolog.INFO)

	// Define a route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/product/:id", handlers.GetProductByIdHandler)

	e.GET("/product", handlers.GetProductsHandler)

	e.POST("/product", handlers.AddProductHandler)

	// Start server
	go func() {
		address := fmt.Sprintf(":%d", config.Port)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
