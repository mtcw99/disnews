package core

import (
	"html/template"
	"net/http"
)

// Paths information for HTML render functionalities
type RenderInfo struct {
	PathTemplates string
	PathStatic    string
}

// Global RenderInfo variable
var Info = RenderInfo{
	PathTemplates: "./templates/",
	PathStatic:    "./static/"}

// Renders the given template file (templatePath), given base.html exists and used
func RenderTemplate(w http.ResponseWriter, templatePath string) {
	t, err := template.ParseFiles(Info.PathTemplates+templatePath,
		Info.PathTemplates+"base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
