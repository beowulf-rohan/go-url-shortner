package controller

import (
	"log"

	"github.com/beowulf-rohan/go-url-shortner/api/services"
	gin "github.com/gin-gonic/gin"
)

func Resolve(c *gin.Context) {
	shortUrl := c.Param("shortURL")
	ip := c.ClientIP()

	resolvedDoc, code, err := services.Resolve(shortUrl, ip)
	if err != nil {
		c.JSON(code, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("Resloved URL:", resolvedDoc)
	
	if code == 200 && resolvedDoc.URL != "" {
        c.Redirect(302, resolvedDoc.URL)
        return
    }

    c.JSON(code, resolvedDoc)
}
