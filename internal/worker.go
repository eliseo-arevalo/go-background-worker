package internal

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "time"
)

func CallAPI(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        log.Printf("❌ Error creating request for %s: %v", url, err)
        return err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("❌ Error calling %s: %v", url, err)
        return err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("❌ Error reading response from %s: %v", url, err)
        return err
    }

    if resp.StatusCode >= 400 {
        log.Printf("❌ Error from %s - Status: %d, Response: %s", url, resp.StatusCode, string(body))
        return fmt.Errorf("status code: %d", resp.StatusCode)
    }

    log.Printf("✅ Success from %s - Response: %s", url, string(body))
    return nil
}

func StartWorker(ctx context.Context, cfg APIConfig) error {
    duration, err := time.ParseDuration(cfg.Interval)
    if err != nil {
        return fmt.Errorf("invalid interval %s: %w", cfg.Interval, err)
    }

    ticker := time.NewTicker(duration)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := CallAPI(ctx, cfg.URL); err != nil {
                log.Printf("❌ Error: %v", err)
            }
        case <-ctx.Done():
            return nil
        }
    }
}