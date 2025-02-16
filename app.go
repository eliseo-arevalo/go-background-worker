package main

import (
   "context"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "os"
   "os/signal" 
   "syscall"
   "time"

   "github.com/joho/godotenv"
)

type APIConfig struct {
   URL      string `json:"url"`
   Interval string `json:"interval"`
}

type APIError struct {
   URL string
   Err error
}

func (e *APIError) Error() string {
   return fmt.Sprintf("API %s error: %v", e.URL, e.Err)
}

func callAPI(ctx context.Context, url string) error {
   req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
   if err != nil {
       log.Printf("❌ Error creating request for %s: %v", url, err)
       return &APIError{URL: url, Err: err}
   }

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
       log.Printf("❌ Error calling %s: %v", url, err)
       return &APIError{URL: url, Err: err}
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
       log.Printf("❌ Error reading response from %s: %v", url, err)
       return &APIError{URL: url, Err: err}
   }

   if resp.StatusCode >= 400 {
       log.Printf("❌ Error from %s - Status: %d, Response: %s", url, resp.StatusCode, string(body))
       return &APIError{URL: url, Err: fmt.Errorf("status code: %d", resp.StatusCode)}
   }

   log.Printf("✅ Success from %s - Response: %s", url, string(body))
   return nil
}

func startAPICaller(ctx context.Context, api APIConfig) error {
   duration, err := time.ParseDuration(api.Interval)
   if err != nil {
       return fmt.Errorf("invalid interval %s: %w", api.Interval, err)
   }

   ticker := time.NewTicker(duration)
   defer ticker.Stop()

   for {
       select {
       case <-ticker.C:
           if err := callAPI(ctx, api.URL); err != nil {
               log.Printf("❌ Error: %v", err)
           }
       case <-ctx.Done():
           return nil
       }
   }
}

func main() {
   logFile, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
   if err != nil {
       fmt.Printf("❌ Error opening log file: %v\n", err)
       return
   }
   defer logFile.Close()

   log.SetOutput(io.MultiWriter(os.Stdout, logFile))

   if err := godotenv.Load(); err != nil {
       log.Println("⚠️ .env file not found, using existing environment variables")
   }

   apisEnv := os.Getenv("APIS")
   if apisEnv == "" {
       log.Fatal("❌ APIS environment variable not found")
   }

   var apiConfigs []APIConfig
   if err := json.Unmarshal([]byte(apisEnv), &apiConfigs); err != nil {
       log.Fatalf("❌ Error parsing APIS: %v", err)
   }

   ctx, cancel := context.WithCancel(context.Background())
   defer cancel()

   for _, api := range apiConfigs {
       go func(api APIConfig) {
           if err := startAPICaller(ctx, api); err != nil {
               log.Printf("❌ Worker error: %v", err)
           }
       }(api)
   }

   log.Println("✅ Background workers started")

   sigChan := make(chan os.Signal, 1)
   signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
   <-sigChan

   log.Println("⚠️ Shutting down workers...")
   cancel()
   time.Sleep(time.Second)
}