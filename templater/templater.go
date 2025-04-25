package templater

import (
	"context"
	"html/template"
	"strings"
)

type HtmlTemplater struct{}

func NewHtmlTemplater() *HtmlTemplater {
	return &HtmlTemplater{}
}

func (*HtmlTemplater) FillTemplate(ctx context.Context, tmpl string, tmplData any) (string, error) {
	templater, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return "", err
	}

	writer := new(strings.Builder)

	if err := templater.Execute(writer, tmplData); err != nil {
		return "", err
	}

	return writer.String(), nil
}
