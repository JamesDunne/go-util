package web

import "strings"

// matches path against "/a/b" routes or "/a/b/*" routes and returns "*" portion or "".
func MatchSimpleRoute(path, route string) (remainder string, ok bool) {
	if path == route {
		return "", true
	}

	if strings.HasPrefix(path, route+"/") {
		return path[len(route)+1:], true
	}

	return "", false
}
