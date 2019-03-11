package templates

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"text/template"
)

func RenderTemplateToBuffer(config interface{}, buffer bytes.Buffer, tpl string) bytes.Buffer {
	tmpl := template.Must(template.New("main.tf").Parse(tpl))
	err := tmpl.Execute(&buffer, config)

	if err != nil {
		logrus.Fatal("Failed writing to Buffer: ", err)
	}

	return buffer
}
