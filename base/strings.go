package base

import "strings"

// Join paths components, preserving leading and trailing '/' chars:
func PathJoin(a, b string) string {
	if strings.HasSuffix(a, "/") && strings.HasPrefix(b, "/") {
		return a + b[1:]
	} else if strings.HasSuffix(a, "/") {
		return a + b
	} else if strings.HasPrefix(b, "/") {
		return a + b
	} else {
		return a + "/" + b
	}
}

func RemovePrefix(s, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return s
	}
	return s[len(prefix):]
}

func RemoveSuffix(s, suffix string) string {
	if !strings.HasSuffix(s, suffix) {
		return s
	}
	return s[0 : len(s)-len(suffix)]
}
