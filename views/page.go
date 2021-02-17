package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/resources"
	"github.com/forewing/goldennum/views/i18n"
	"github.com/gin-gonic/gin"
)

const (
	templateIndex   = "index.html"
	templateAdmin   = "admin.html"
	templateBaseURL = "base_url"
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

// MustLoadTemplate load template
func MustLoadTemplate() *template.Template {
	templates := resources.GetTemplates()

	// load general templates
	t, err := template.ParseFS(templates, "*.html")
	if err != nil {
		panic(err)
	}

	// generate BaseURL template
	t, err = t.New(templateBaseURL).Parse(
		fmt.Sprintf("{{ define \"%v\" }}%v{{ end }}", templateBaseURL, config.Load().BaseURL),
	)
	if err != nil {
		panic(err)
	}
	return t
}
