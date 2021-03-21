package core

import (
	"html/template"
	"net/http"

	"github.com/mtcw99/disnews/sessions"
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
func RenderTemplate(w http.ResponseWriter, templatePath string,
		sessionInfo *sessions.SessionInfo, data interface{}) error {
	t, err := template.ParseFiles(Info.PathTemplates+templatePath,
		Info.PathTemplates+"base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
