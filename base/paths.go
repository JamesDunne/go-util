package base

import "path/filepath"

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
