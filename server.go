package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kengathua/book-inventory-system/pkg/apis"
	gormClient "github.com/Kengathua/book-inventory-system/pkg/infrastructure/database/gorm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	dbClient = &gormClient.DBClient{}
)

func main() {
	dbClient = gormClient.NewPGDBClient()
	app := fiber.New(fiber.Config{
		ProxyHeader: "trust",
	})
	app.Use(cors.New())

	app.Static("/register", "./pkg/templates/register.html")
	app.Static("/login", "./pkg/templates/login.html")
	app.Static("/books", "./pkg/templates/books.html")

	apisURL := app.Group("/apis")
	apis.RegisterAPIsRoutes(apisURL, dbClient.DB)
	apis.RegisterAuthRoutes(apisURL, dbClient.DB)

	args := os.Args[1:]
	portNumber := "8000" // default port number
	if len(args) > 0 {
		portNumber = args[0]
	}

	port := fmt.Sprintf(":%s", portNumber)

	// Start the server in a separate goroutine.
	go func() {
		if err := app.Listen(port); err != nil {
			if isPortInUseError(err) {
				log.Printf("Port %s is already in use. Exiting gracefully...\n", portNumber)
			} else {
				log.Printf("Server error: %v\n", err)
			}
		}
	}()

	// Serve the app and wait for an interrupt signal (e.g., SIGINT or SIGTERM).
	waitForInterrupt(app)
}

func isPortInUseError(err error) bool {
	if opError, ok := err.(*net.OpError); ok {
		if sysErr, ok := opError.Err.(*os.SyscallError); ok {
			if sysErr.Err == syscall.EADDRINUSE {
				return true
			}
		}
	}

	return false
}

func waitForInterrupt(app *fiber.App) {
	// Creating a channel to capture the interrupt signal.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt) // Capture SIGINT (Ctrl+C) signals.

	// Wait for the interrupt signal.
	<-interruptChan

	log.Println("Received interrupt signal (Ctrl+C). Shutting down gracefully...")

	// Creating a context with a timeout for the shutdown process.
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Adjustible timeout.
	defer cancel()

	// Using server.Shutdown to gracefully shut down the Fiber application.
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during clean shutdown: %s\n", err)
		os.Exit(1) // Exit with an error code
	} else {
		log.Println("Shut down gracefully.")
	}
}
