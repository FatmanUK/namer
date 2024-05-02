package main

import (
	"strings"
)

func StringFromStringArray(strs []string) string {
	return "[" + strings.Join(strs, " ") + "]"
}
