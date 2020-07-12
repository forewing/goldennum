package views

import (
	"github.com/gin-gonic/gin"
)

func AdminIndex(c *gin.Context) {
	c.String(200, "fuck")
}
