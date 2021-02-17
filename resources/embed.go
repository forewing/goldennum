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
)

var (
	//go:embed statics/*
	staticsEmbed embed.FS
	statics      fs.FS

	//go:embed templates/*
	templatesEmbed embed.FS
	templates      fs.FS
)

func init() {
	statics = mustStripPrefix(staticsEmbed, staticsPath)
	templates = mustStripPrefix(templatesEmbed, templatesPath)
}

// SetLiveReload make resources load from disk
// if avaliable, instead of from embedded files
func SetLiveReload() {
	setLiveReload(&statics, staticsPath)
	setLiveReload(&templates, templatesPath)
}

// GetStatics return statics
func GetStatics() fs.FS {
	return statics
}

// GetTemplates return templates
func GetTemplates() fs.FS {
	return templates
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
