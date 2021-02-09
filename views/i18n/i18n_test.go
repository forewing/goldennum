package i18n

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func TestCombineLanguageConfig(t *testing.T) {
	base := LanguageConfig{
		"a": "1",
		"b": "2",
	}
	override := LanguageConfig{
		"b": "3",
		"c": "4",
	}
	expect := LanguageConfig{
		"a": "1",
		"b": "3",
	}
	ret := combineLanguageConfig(base, override)
	if !reflect.DeepEqual(ret, expect) {
		t.Errorf("not equal: %v, %v", ret, expect)
	}
}

func TestGetI18nData(t *testing.T) {
	Load()
	languageHeaderMap := map[language.Tag]string{
		language.English:            "en",
		language.SimplifiedChinese:  "zh-Hans",
		language.TraditionalChinese: "zh-Hant",
		language.Chinese:            "zh",
	}
	for i, l := range i18nLanguages {
		c := &gin.Context{
			Request: &http.Request{
				Header: http.Header{},
			},
		}
		c.Request.Header.Add(acceptLanguageHeaderKey, languageHeaderMap[l])
		t.Log(i)
		data := GetI18nData(c)
		t.Log(i18nCachedData, i, languageHeaderMap[l])
		if !reflect.DeepEqual(data, i18nCachedData[i]) {
			t.Errorf("%v not equal: %v %v", acceptLanguageHeaderKey, i, i18nConfigPath[i])
		}
	}
}
