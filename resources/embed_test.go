package resources

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"testing"
)

func TestEmbedStatics(t *testing.T) {
	err := testEmbed(GetStatics(), staticsPath)
	if err != nil {
		t.Error(err)
	}
}

func TestEmbedTemplates(t *testing.T) {
	err := testEmbed(GetTemplates(), templatesPath)
	if err != nil {
		t.Error(err)
	}
}

func TestEmbedI18n(t *testing.T) {
	err := testEmbed(GetI18n(), i18nPath)
	if err != nil {
		t.Error(err)
	}
}

func testEmbed(vfs fs.FS, root string) error {
	dir, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, file := range dir {
		if !file.Type().IsRegular() {
			continue
		}
		name := file.Name()
		vf, err := fs.ReadFile(vfs, name)
		if err != nil {
			return err
		}
		rf, err := os.ReadFile(path.Join(root, name))
		if err != nil {
			return err
		}
		if res := bytes.Compare(vf, rf); res != 0 {
			return fmt.Errorf("%v/%v not match", root, name)
		}
	}
	return nil
}
