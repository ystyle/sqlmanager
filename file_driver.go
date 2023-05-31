package sqlmanager

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MarkdownDriver struct {
	dir string
}

type item struct {
	Type    string
	Content string
}

func NewMarkdownDriver() *MarkdownDriver {
	return &MarkdownDriver{
		dir: "./sql",
	}
}

func NewMarkdownDriverWithDir(dir string) *MarkdownDriver {
	return &MarkdownDriver{
		dir: dir,
	}
}
func (mdd *MarkdownDriver) DriverName() string {
	return "Markdown"
}

func (mdd *MarkdownDriver) Load() ([]SqlTemple, error) {
	var sqls []SqlTemple
	err := filepath.Walk(mdd.dir, func(subpath string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		ext := path.Ext(subpath)
		if ext == ".md" || ext == ".markdown" {
			s, err := mdd.parseMarkdown(subpath)
			if err != nil {
				return err
			}
			if len(s) != 0 {
				sqls = append(sqls, s...)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sqls, nil
}

func (mdd *MarkdownDriver) parseMarkdown(filename string) ([]SqlTemple, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("sqlmanager - ERROR: %s loading failed...\n", filename)
		return nil, err
	}
	return parseMarkdown(buf, mdd.getName(filename))
}

func (mdd *MarkdownDriver) getName(filename string) string {
	ext := path.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	base = strings.TrimPrefix(base, mdd.dir)
	base = strings.TrimPrefix(base, "/")
	return base
}
