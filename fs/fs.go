package fs

import (
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CanonicalPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		panic(err)
	}
	return abs
}

func GetMimeType(filename string) string {
	return mime.TypeByExtension(strings.ToLower(path.Ext(filename)))
}

func ExtractNames(fis []os.FileInfo) []string {
	names := make([]string, len(fis), len(fis))
	for i := range fis {
		names[i] = fis[i].Name()
	}
	return names
}
