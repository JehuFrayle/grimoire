package utils

import (
	"strings"
)

func ExtractID(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	id := strings.TrimPrefix(path, prefix)
	id = strings.Trim(id, "/")
	return id
}
