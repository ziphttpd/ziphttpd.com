package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var templates *Template

// Template is a custom html/template renderer for Echo framework
type Template struct {
	templates *template.Template
}

// Render は echo.Renderer の実装
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// レンダリング実行
func templateRender(code int, name string, data interface{}, c echo.Context) error {
	if templates == nil {
		templates = &Template{
			templates: template.Must(template.ParseGlob("template/*.html")),
		}
		c.Echo().Renderer = templates
	}
	return c.Render(code, name, data)
}
