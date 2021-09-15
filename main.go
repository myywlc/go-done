package main

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
    r := gin.Default()
    
    // Ping test
    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "ping111")
    })
    
    // Get user value
    r.GET("/user/:name", func(c *gin.Context) {
        user := c.Params.ByName("name")
        value, ok := db[user]
        if ok {
            c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
        } else {
            c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
        }
    })
    
    authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
        "foo":  "bar", // user:foo password:bar
        "manu": "123", // user:manu password:123
    }))
    
    authorized.POST("admin", func(c *gin.Context) {
        user := c.MustGet(gin.AuthUserKey).(string)
        
        // Parse JSON
        var json struct {
            Value string `json:"value" binding:"required"`
        }
        
        if c.Bind(&json) == nil {
            db[user] = json.Value
            c.JSON(http.StatusOK, gin.H{"status": "ok"})
        }
    })
    
    return r
}

func main() {
    r := setupRouter()
    // Listen and Server in 0.0.0.0:8080
    err := r.Run("127.0.0.1:3000")
    if err != nil {
        return
    }
}
