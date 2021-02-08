package i18n

import (
	"io"

	"github.com/forewing/goldennum/resources"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// LanguageConfig configs for all pages
type LanguageConfig gin.H

var (
	i18nLanguages = []language.Tag{
		language.English,
		language.SimplifiedChinese,
		language.TraditionalChinese,
		language.Chinese,
	}
	i18nMatcher = language.NewMatcher(i18nLanguages)

	i18nConfigPath = []string{
		"en.yml", // The first is the fallback config
		"zh-Hans.yml",
		"zh-Hant.yml",
		"zh-Hans.yml",
	}

	i18nCachedData = []LanguageConfig{}
)

// Load all i18n configs
func Load() {
	for i, path := range i18nConfigPath {
		i18nLoadData := parseI18nConfig(path)
		if i > 0 {
			i18nLoadData = combineLanguageConfig(i18nCachedData[0], i18nLoadData)
		}
		i18nCachedData = append(i18nCachedData, i18nLoadData)
		zap.S().Infof("[i18n] Load %v success", path)
	}
	zap.S().Debug(i18nCachedData)
}

func combineLanguageConfig(base LanguageConfig, override LanguageConfig) LanguageConfig {
	result := LanguageConfig{}
	for k, v := range base {
		result[k] = v
		if v2, ok := override[k]; ok {
			result[k] = v2
		}
	}
	return result
}

func parseI18nConfig(path string) LanguageConfig {
	file, err := resources.GetI18n().Open(path)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	i18nData := LanguageConfig{}
	err = yaml.Unmarshal(data, &i18nData)
	if err != nil {
		panic(err)
	}
	return i18nData
}

// GetI18nData returns target language data
func GetI18nData(c *gin.Context) LanguageConfig {
	accept := c.GetHeader("Accept-Language")
	_, i := language.MatchStrings(i18nMatcher, accept)
	return i18nCachedData[i]
}
