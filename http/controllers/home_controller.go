package controllers

import (
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Home renders the home page using the home.html template.
func Home(c *gin.Context) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":     "Welcome to RapidGo",
		"version":   "0.2.0",
		"goVersion": runtime.Version(),
		"env":       env,
	})
}
