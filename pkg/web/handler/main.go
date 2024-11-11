package handler

import (
	"html/template"
	"log"
	"net/http"
)

var basePath = "./pkg/web/views/"

func SetBasePath(path string) {
	basePath = path
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("landing.html.tmpl")
	tmpl, err := tmpl.ParseFiles(basePath + "landing.html.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, map[string]any{
		"title": "This is a title",
	})
	if err != nil {
		log.Println(err)
	}
}
