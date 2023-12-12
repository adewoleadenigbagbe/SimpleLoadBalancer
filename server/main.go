package main

import "fmt"

func main() {

	app.Echo.Logger.SetLevel(log.INFO)

	// Define a route
	app.Echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Start server
	go func() {
		if err := app.Echo.Start(":8185"); err != nil && err != http.ErrServerClosed {
			app.Echo.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Echo.Shutdown(ctx); err != nil {
		app.Echo.Logger.Fatal(err)
	}
}
