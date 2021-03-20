package views

import (
	"io"
	"testing"

	"github.com/forewing/goldennum/version"
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

func TestAddURLVersionQuery(t *testing.T) {
	cases := [][4]string{
		{"/goldennum", "statics/a.js", "", "/goldennum/statics/a.js"},
		{"/goldennum", "/statics/a.js", "", "/goldennum/statics/a.js"},
		{"goldennum", "statics/a.js", "", "goldennum/statics/a.js"},
		{"goldennum", "/statics/a.js", "", "goldennum/statics/a.js"},
		{"", "/statics/a.js", version.HashDefault, "/statics/a.js"},
		{"", "/statics/a.js", "", "/statics/a.js"},
		{"", "/statics/a.js", "1234", "/statics/a.js"},
		{"", "/statics/a.js", "1234567", "/statics/a.js?v=1234567"},
		{"", "/statics/a.js", "123456789", "/statics/a.js?v=1234567"},
		{"", "/statics/a.js", "1/234567", "/statics/a.js?v=1%2F23456"},
		{"", "/statics/a.js?a=69420", "1234567", "/statics/a.js?a=69420&v=1234567"},
		{"", "/statics/a.js?z=69420", "1234567", "/statics/a.js?v=1234567&z=69420"},
	}
	for i, c := range cases {
		r := generateStaticURL(c[0], c[1], c[2])
		if r != c[3] {
			t.Errorf("case %v: base=%v, origin=%v, Hash=%v, got %v, expect %v", i, c[0], c[1], c[2], r, c[3])
		}
	}
}
