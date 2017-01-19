package template

import (
	"html/template"
	"strings"
)

func TemplateHelpers() template.FuncMap {
	return template.FuncMap{
		"capitalize": Capitalize,
	}
}

// Capitalize capitalizes the given string.
func Capitalize(str string) string {
	if len(str) <= 0 {
		return ""
	}

	if len(str) == 1 {
		return strings.ToUpper(str)
	}

	return strings.ToUpper(str[:1]) + strings.ToLower(str[1:])
}
