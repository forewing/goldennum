package views

import (
	"net/http"

	"github.com/forewing/goldennum/views/i18n"
	"github.com/gin-gonic/gin"
)

const (
	templateIndex = "index.html"
	templateAdmin = "admin.html"
)

// PageIndex render index
func PageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateIndex, gin.H{
		"i18n": i18n.GetI18nData(c),
	})
}

// AdminIndex return admin page
func AdminIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateAdmin, gin.H{
		"i18n": i18n.GetI18nData(c),
	})
}
