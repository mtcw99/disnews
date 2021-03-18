package core

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

type RenderInfo struct {
	Path string
}

var Info = RenderInfo{
	Path: "./templates/"}

func RenderTemplate(w http.ResponseWriter, templatePath string) {
	t, err := template.ParseFiles(Info.Path + templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
