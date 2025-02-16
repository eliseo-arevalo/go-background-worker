package internal

import (
    "os"

    "github.com/gin-gonic/gin"
)

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