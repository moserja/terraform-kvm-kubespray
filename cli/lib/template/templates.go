package template

import (
	"bytes"
	"os"
	"strings"
	"text/template"
	"unicode"
)

type Template interface {
	Name() string
	Functions() map[string]any
}

type TextTemplate interface {
	Template
	Template() string
}

// Populate populates the template and returns it as a string.
func populate(t Template, content string) (string, error) {
	tpl := template.New(t.Name())
	tpl = tpl.Funcs(BuiltInFuncs())
	tpl = tpl.Funcs(template.FuncMap(t.Functions()))

	tpl, err := tpl.Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tpl.Execute(&buf, t)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Populate populates the template and returns it as a string.
func Populate(t TextTemplate) (string, error) {
	return populate(t, t.Template())
}

// PopulateFrom reads the template from a given path and returns
// the populated template as a string.
func PopulateFrom(t Template, path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return populate(t, string(content))
}

// Write populates the template and writes it to the given path.
func Write(t TextTemplate, path string) error {
	res, err := Populate(t)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(res), os.ModePerm)
}

// WriteFrom reads a template from source path, populates the template
// and writes it to the destination path.
func WriteFrom(t Template, srcPath string, dstPath string) error {
	res, err := PopulateFrom(t, srcPath)
	if err != nil {
		return err
	}

	return os.WriteFile(dstPath, []byte(res), os.ModePerm)
}

// TrimTemplate trims alls leading and trailing spaces from each line
// and replace tabs with double space.
func TrimTemplate(s string) string {
	var l = strings.Split(s, "\n")

	for i, s := range l {
		l[i] = strings.ReplaceAll(s, "\t", "  ")
	}

	return strings.Join(trimLines(l), "\n")
}

// trimLines removes maximum equal prefix spaces from each line.
// In addition, it removes empty leading and trailing lines
func trimLines(l []string) []string {
	l = trimEmptyLines(l)
	indent := -1

	// Evaluate max leading spaces
	for _, s := range l {
		if s == "" {
			continue
		}

		ls := leadingSpaces(s)
		if indent == -1 || ls < indent {
			indent = ls
		}
	}

	// Remove max leading spaces
	for i := range l {
		l[i] = strings.Replace(l[i], " ", "", indent)
		l[i] = strings.TrimRight(l[i], " ")
	}

	return trimEmptyLines(l)
}

// trimEmptyLines removes leading and trailing empty lines.
func trimEmptyLines(l []string) []string {
	l = trimLeadingEmptyLines(l)
	l = trimTrailingEmptyLines(l)
	return l
}

// trimLeadingEmptyLines removes leading empty lines.
func trimLeadingEmptyLines(l []string) []string {
	if len(l) == 0 || !empty(l[0]) {
		return l
	}

	return trimLeadingEmptyLines(l[1:])
}

// trimTrailingEmptyLines removes empty trailing lines.
func trimTrailingEmptyLines(l []string) []string {
	last := len(l) - 1

	if last < 0 || !empty(l[last]) {
		return l
	}

	return trimTrailingEmptyLines(l[:last])
}

// leadingSpaces returns count of leading spaces in a string.
func leadingSpaces(s string) int {
	count := 0

	for _, r := range s {
		if !unicode.IsSpace(r) {
			return count
		}
		count++
	}

	return count
}

// empty returns true if string contains only spaces, newlines and tabs.
func empty(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}
