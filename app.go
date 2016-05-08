package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/commonlog"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Use(static.Serve("/js", static.LocalFile("public/js", true)))

	r.Use(commonlog.New())

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	port := os.Getenv("PORT")
	if port == "" {
		r.Run(":5000")
	} else {
		r.Run(":" + port)
	}
}
