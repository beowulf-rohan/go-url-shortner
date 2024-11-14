package controller

import (
	"log"

	"github.com/beowulf-rohan/go-url-shortner/api/services"
	gin "github.com/gin-gonic/gin"
)

func Resolve(c *gin.Context) {
	url := c.Param("url")
	log.Println("URL received for resolution:", url)

	resolvedURL, code, err := services.Resolve(url)
	if err != nil {
		c.JSON(code, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("Resloved URL:", resolvedURL)

	c.Redirect(code, resolvedURL)
}
