package web

import (
	"bytes"
	"embed"
	"io"
	"net/url"
	"text/template"
)

type TemplateData struct {
	Data interface{}
}

//go:embed templates
var templates embed.FS

var funcMap = template.FuncMap{
	"urljoin": func(elems ...string) string {
		r, _ := url.JoinPath(elems[0], elems[1:]...)
		return r
	},
}

func CreateTemplateWithBase(templateName string) (*template.Template, error) {

	return template.New(templateName).Funcs(funcMap).ParseFS(
		templates,
		"templates/base.tmpl",
		"templates/"+templateName+".tmpl",
	)
}

func RenderTemplateWithBase(w io.Writer, templateName string, data interface{}) error {

	t, err := CreateTemplateWithBase(templateName)

	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, "base", TemplateData{
		Data: data,
	})

	return err

}

func RenderTemplateToString(templateName string, data interface{}) (string, error) {
	tmplStr, _ := templates.ReadFile("templates/" + templateName + ".tmpl")

	t, err := template.New("templates/" + templateName + ".tmpl").Funcs(funcMap).Parse(string(tmplStr))

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	err = t.Execute(&output, data)
	return output.String(), err
}
