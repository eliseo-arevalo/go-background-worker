package internal

import (
    "os"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey != os.Getenv("API_KEY") {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func GetLogs(c *gin.Context) {
    content, err := os.ReadFile("api.log")
    if err != nil {
        c.JSON(500, gin.H{"error": "Could not read logs"})
        return
    }

    c.JSON(200, gin.H{
        "logs": string(content),
    })
}