package compiler

import (
	"fmt"
	"strings"
)

const (
	templateIdentifier = '%'
)

// TemplateValues are the values to be interpolated in a template string.
type TemplateValues map[string]*string

// ParseTemplate parses a Javascript template string, interpolating values that
// are surrounded by % characters. Returns the interpolated string, which
// should be valid Javascript code.
func ParseTemplate(template string, values TemplateValues) string {
	var interpolating bool
	var optKeys []string
	var optBuff []string

	// Remove linebreak on first line of multiline string.
	if template[0] == '\n' {
		template = template[1:]
	}

	for _, r := range template {
		// Check if it is an the interpolation character and begin interpolating on the next rune.
		if r == templateIdentifier {
			if !interpolating {
				optBuff = nil
			} else {
				optKeys = append(optKeys, strings.Join(optBuff, ""))
			}
			interpolating = !interpolating
			continue
		}

		if interpolating {
			optBuff = append(optBuff, string(r))
		}
	}

	for _, k := range optKeys {
		v := values[k]
		if v == nil {
			panic("no replacement value found for key: " + k)
		}
		template = strings.ReplaceAll(template, fmt.Sprintf("%c%s%c", templateIdentifier, k, templateIdentifier), *v)
	}

	return template
}
