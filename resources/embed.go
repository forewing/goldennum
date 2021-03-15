package resources

import (
	"embed"
	"io/fs"
	"os"
	"path"

	"go.uber.org/zap"
)

const (
	packagePath   = "resources"
	staticsPath   = "statics"
	templatesPath = "templates"
	i18nPath      = "i18n"
)

var (
	//go:embed statics/*
	staticsEmbed embed.FS
	statics      = mustStripPrefix(staticsEmbed, staticsPath)

	//go:embed templates/*
	templatesEmbed embed.FS
	templates      = mustStripPrefix(templatesEmbed, templatesPath)

	//go:embed i18n/*
	i18nEmbed embed.FS
	i18n      = mustStripPrefix(i18nEmbed, i18nPath)
)

// SetLiveReload make resources load from disk
// if available, instead of from embedded files
func SetLiveReload() {
	setLiveReload(&statics, staticsPath)
	setLiveReload(&templates, templatesPath)
	setLiveReload(&i18n, i18nPath)
}

// GetStatics return statics
func GetStatics() fs.FS {
	return statics
}

// GetTemplates return templates
func GetTemplates() fs.FS {
	return templates
}

// GetI18n return i18n configs
func GetI18n() fs.FS {
	return i18n
}

func mustStripPrefix(sfs fs.FS, prefix string) fs.FS {
	dfs, err := fs.Sub(sfs, prefix)
	if err != nil {
		panic(err)
	}
	return dfs
}

func setLiveReload(target *fs.FS, name string) {
	realPath := path.Join(packagePath, name)
	if _, err := os.Stat(realPath); err != nil {
		if os.IsNotExist(err) {
			// Cannot live reload
			return
		}
	}
	zap.S().Infof("live reload %v/*", realPath)
	*target = os.DirFS(realPath)
}
