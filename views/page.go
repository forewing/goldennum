package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	templateIndex = 1
)

var (
	// Templates name lookup table
	Templates map[int]string = map[int]string{
		templateIndex: "index.html",
	}
)

// PageIndex render index
func PageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, Templates[templateIndex], gin.H{})
}
