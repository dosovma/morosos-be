package templater

import (
	"html/template"
	"strings"
)

type HTMLTemplater struct{}

func NewHTMLTemplater() *HTMLTemplater {
	return &HTMLTemplater{}
}

func (*HTMLTemplater) FillTemplate(tmpl string, tmplData any) (string, error) {
	templater, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return "", err
	}

	writer := new(strings.Builder)

	if err = templater.Execute(writer, tmplData); err != nil {
		return "", err
	}

	return writer.String(), nil
}
