package main

import (
    "context"
    "io"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "bworker/internal"
)

func setupLogging() error {
    logFile, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return err
    }
    log.SetOutput(io.MultiWriter(os.Stdout, logFile))
    return nil
}

func main() {
    if err := setupLogging(); err != nil {
        log.Fatalf("❌ Error setting up logging: %v", err)
    }

    configs, err := internal.LoadConfig()
    if err != nil {
        log.Fatalf("❌ Error loading config: %v", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Start workers
    for _, cfg := range configs {
        go func(cfg internal.APIConfig) {
            if err := internal.StartWorker(ctx, cfg); err != nil {
                log.Printf("❌ Worker error: %v", err)
            }
        }(cfg)
    }

    // Setup and start Gin server
    router := gin.Default()
    router.GET("/logs", internal.AuthMiddleware(), internal.GetLogs)
    
    go func() {
        if err := router.Run(":8080"); err != nil {
            log.Printf("❌ Server error: %v", err)
        }
    }()

    log.Println("✅ Workers and API server started")

    // Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    log.Println("⚠️ Shutting down...")
    cancel()
    time.Sleep(time.Second)
}