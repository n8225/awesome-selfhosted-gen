package parse

import (
	"strings"
)

func getHeader(i int, l string) string {
	if i == 0 && !strings.HasPrefix(l, "-") {
		return l
	}
	return ""
}
