package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	templateIndex = "index.html"
	templateAdmin = "admin.html"
)

// PageIndex render index
func PageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateIndex, gin.H{})
}

// AdminIndex return admin page
func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateAdmin, gin.H{})
}
