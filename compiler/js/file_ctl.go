package js

import (
	"fmt"
	"os"
	"strings"
)

const (
	templateLiteralChar = "%"
)

// TemplateValues are the values to be interpolated in a template string.
type TemplateValues map[string]string

// ParseTemplate parses a Javascript template string, interpolating values that
// are surrounded by % characters. Returns the interpolated string, which
// should be valid Javascript code.
func ParseTemplate(template string, values TemplateValues) string {
	var interpolating bool
	var optKeys []string
	var optBuff []string

	for _, r := range template {
		c := string(r)

		// Check if it is an the interpolation character and begin interpolating on the next character.
		if c == templateLiteralChar {
			if !interpolating {
				optBuff = nil
			} else {
				optKeys = append(optKeys, strings.Join(optBuff, ""))
			}
			interpolating = !interpolating
			continue
		}

		if interpolating {
			optBuff = append(optBuff, c)
		}
	}

	for _, k := range optKeys {
		v := values[k]
		if v == "" {
			panic("no replacement value found for key: " + k)
		}
		template = strings.ReplaceAll(template, fmt.Sprintf("%s%s%s", templateLiteralChar, k, templateLiteralChar), v)
	}

	return template
}

// WriteToFile writes a byte slice to a file.
func (f FileCtl) WriteToFile(name string, contents []byte) {
	file, err := os.OpenFile(f.concatToPath(name), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := file.Write(contents); err != nil {
		panic(err)
	}
}

// InitDir initializes the working directory.
func (f FileCtl) InitDir() {
	if err := os.MkdirAll(f.Directory, 0777); err != nil {
		panic(err)
	}
}

// concatToPath joins a file name to the working directory.
func (f FileCtl) concatToPath(name string) string {
	return f.Directory + "/" + name
}

func NewFileCtl(directory string, joins ...*FileCtl) *FileCtl {
	ctl := &FileCtl{Directory: directory}
	for _, j := range joins {
		ctl.Directory = fmt.Sprintf("%s/%s", ctl.Directory, j.Directory)
	}
	return ctl
}

// FileCtl is a utility struct for writing Javascript files.
type FileCtl struct {
	Directory string
}

// FileCtlOption is an option for configuring a FileCtl.
type FileCtlOption func(*FileCtl)
