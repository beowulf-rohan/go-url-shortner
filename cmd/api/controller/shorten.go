package controller

import (
	"github.com/asaskevich/govalidator"
	"github.com/beowulf-rohan/go-url-shortner/api/services"
	"github.com/beowulf-rohan/go-url-shortner/model"
	"github.com/beowulf-rohan/go-url-shortner/utils"
	"github.com/gin-gonic/gin"
)

func Shorten(c *gin.Context) {
	request := model.Request{}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"error": "error while unmarshalling input",
		})
		return
	}

	if !govalidator.IsURL(request.URL) {
		c.JSON(400, gin.H{
			"error": "request URL is not valid",
		})
		return
	}

	if !utils.CheckDomainError(request.URL) {
		c.JSON(503, gin.H{
			"error": "domain not valid",
		})
		return
	}

	request.URL = utils.EnforceHttp(request.URL)

	shortenedURLResponse, code, err := services.Shorten(&request, c.ClientIP())
	if err != nil {
		c.JSON(code, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(code, shortenedURLResponse)
}
