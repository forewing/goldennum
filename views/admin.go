package views

import (
	"github.com/gin-gonic/gin"
)

func AdminIndex(c *gin.Context) {
	c.JSON(200, "fuck")
}
