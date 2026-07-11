package main

import (
    "io"
    "net/http"
    "strings"
    "fmt"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "resumeparser/extractor"
    "resumeparser/parser"
)

func main() {
     // ── API connection check on startup ──
    fmt.Println("Connecting to NVIDIA API...")
    if err := parser.CheckConnection(); err != nil {
        fmt.Println("Connection failed:", err)
        // don't exit — still start server, just warn
    } else {
        fmt.Println("Connected to NVIDIA API successfully!")
    }

    r := gin.Default()

    
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: false,
    }))

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    r.POST("/parse", func(c *gin.Context) {
        file, header, err := c.Request.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
            return
        }
        defer file.Close()

        data, err := io.ReadAll(file)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
            return
        }

        text, err := extractor.Extract(header.Filename, data)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if strings.TrimSpace(text) == "" {
            c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "no text found in file"})
            return
        }

        resume, err := parser.ParseResume(text)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "success", "data": resume})
    })

    r.Run(":8080")
}