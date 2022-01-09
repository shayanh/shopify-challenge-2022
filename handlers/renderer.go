package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Renderer renders some output related to the given name and data to the given
// http response.
type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
}

// HTMLRenderer renders HTML output using templates.
type HTMLRenderer struct {
	templates *template.Template
}

func NewHTMLRenderer(templatesBaseDir string) *HTMLRenderer {
	templateNames := []string{
		"list.html",
		"edit.html",
	}
	var templateFileNames []string
	for _, tn := range templateNames {
		templateFileNames = append(templateFileNames, filepath.Join(templatesBaseDir, tn))
	}
	return &HTMLRenderer{
		templates: template.Must(template.ParseFiles(templateFileNames...)),
	}
}

func (h *HTMLRenderer) Render(w http.ResponseWriter, tmpl string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
