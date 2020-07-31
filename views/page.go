package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/forewing/goldennum/config"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

const (
	templateHeader  = "header.html"
	templateFooter  = "footer.html"
	templateBaseURL = "base_url.html"
	templateIndex   = "index.html"
)

var (
	templatesBox = packr.New("templates", "../templates")

	templates []string = []string{
		templateHeader,
		templateFooter,
		templateIndex,
	}
)

func generateBaseURL() string {
	c := config.Load()
	return fmt.Sprintf("{{ define \"%v\" }}%v{{ end }}", templateBaseURL, c.BaseURL)
}

// LoadTemplate reutrn templates
func LoadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range templates {
		str, err := templatesBox.FindString(name)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(str))
		if err != nil {
			return nil, err
		}
	}
	return t.New(templateBaseURL).Parse(generateBaseURL())
}

// PageIndex render index
func PageIndex(c *gin.Context) {
	c.HTML(http.StatusOK, templateIndex, gin.H{})
}
