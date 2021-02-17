package views

import (
	"io"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTemplateLoad(t *testing.T) {
	MustLoadTemplate()
}

func TestTemplateExecute(t *testing.T) {
	templates := MustLoadTemplate()
	targets := []string{templateIndex, templateAdmin}
	for _, target := range targets {
		err := templates.ExecuteTemplate(io.Discard, target, gin.H{
			"i18n": gin.H{},
		})
		if err != nil {
			t.Error(err)
		}
	}
}
