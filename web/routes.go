package web

import "strings"

// matches path against "/a/b" exact route.
func MatchExactRoute(path, route string) (ok bool) {
	return path == route
}

// matches path against "/a/b/" exact route.
func MatchExactRouteIgnoreSlash(path, route string) (ok bool) {
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	if strings.HasSuffix(route, "/") {
		route = route[:len(route)-1]
	}
	return path == route
}

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

// matches path against "/a/b{*}" routes and returns "*" portion (including leading '/') or "".
func MatchSimpleRouteRaw(path, route string) (remainder string, ok bool) {
	if path == route {
		return "", true
	}

	if strings.HasPrefix(path, route) {
		return path[len(route):], true
	}

	return "", false
}
