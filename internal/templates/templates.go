package templates

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed all:embeds
var FS embed.FS

type Data struct {
	ProjectName string
	Year        int
}

func Process(templatePath string, data Data) ([]byte, error) {
	// The path needs to be relative to the embedded FS root
	fullPath := "embeds/" + templatePath
	fileBytes, err := FS.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	t, err := template.New(templatePath).Parse(string(fileBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
