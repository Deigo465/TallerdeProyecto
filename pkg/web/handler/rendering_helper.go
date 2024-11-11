package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Response struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

func JSON(w http.ResponseWriter, code int, message string, data interface{}) {
	wrappedData := Response{data, code, message}
	// set the content type first
	w.Header().Set("Content-Type", "application/json")
	// then the code
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(wrappedData); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ExecLayout(w io.Writer, tmplName string, templateVars map[string]any) {
	layout := template.New("base_layout.html.tmpl")

	// check if tmplName contains slashes

	route := tmplName
	i := strings.Index(tmplName, "/")
	if i != -1 {
		tmplName = strings.Split(tmplName, "/")[1]
	}

	tmpl := template.New(tmplName)

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	layout = layout.Funcs(map[string]interface{}{
		"title": func() string { return "title" },
		"body": func() error {
			return tmpl.Execute(w, templateVars)
		},
	})
	layout, err := layout.ParseFiles(basePath + "layouts/base_layout.html.tmpl")
	if err != nil {
		log.Println("error parsing layout: ", err)
		return
	}

	tmpl, err = tmpl.ParseFiles(basePath + route)
	if err != nil {
		fmt.Println("error parsing template: ", err)
		return
	}

	err = layout.Execute(w, templateVars)
	if err != nil {
		log.Println("error executing layout : ", err)
	}
}
