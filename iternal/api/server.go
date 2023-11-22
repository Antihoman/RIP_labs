package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/services", func(c *gin.Context) {
		filter := c.Query("filter")
		data := GetAllServices(filter)
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	r.GET("/service/:index", func(c *gin.Context) {
		index := c.Param("index")
		serviceIndex, err := strconv.Atoi(index)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid index")
			return
		}

		service := GetServiceByIndex(serviceIndex)
		c.HTML(http.StatusOK, "detail.tmpl", service)
	})

	r.Static("/image", "./resources/images")
	r.Static("/css", "./resources/css")

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	log.Println("Server down")
}
