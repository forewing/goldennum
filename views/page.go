package views

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"github.com/forewing/goldennum/config"
	"github.com/forewing/goldennum/resources"
	"github.com/forewing/goldennum/version"
	"github.com/forewing/goldennum/views/i18n"
	"github.com/gin-gonic/gin"
)

const (
	templateIndex   = "index.html"
	templateAdmin   = "admin.html"
	templateBaseURL = "base_url"

	queryHashLength = 7
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

	baseURL := config.Load().BaseURL

	// load general templates
	t, err := template.New("").Funcs(template.FuncMap{
		"generateStaticURL": func(path string) string {
			return generateStaticURL(baseURL, path, version.Hash)
		},
	}).ParseFS(templates, "*.html")

	if err != nil {
		panic(err)
	}

	// generate BaseURL template
	t, err = t.New(templateBaseURL).Parse(
		fmt.Sprintf("{{ define \"%v\" }}%v{{ end }}", templateBaseURL, baseURL),
	)
	if err != nil {
		panic(err)
	}
	return t
}

func generateStaticURL(base string, origin string, hash string) string {
	u, err := url.Parse(origin)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(base, u.Path)

	if hash != version.HashDefault && len(hash) >= queryHashLength {
		q := u.Query()
		q.Set("v", hash[0:queryHashLength])
		u.RawQuery = q.Encode()
	}

	return u.String()
}
