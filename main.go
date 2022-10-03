package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func main(){
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		var filepath = path.Join("views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var data = map[string]interface{}{
			"title" : "Beta Golang Web Project", 
			"name" : "Sembada tech",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("assets"))))

	fmt.Println("server started at localhost: 8000")
	http.ListenAndServe(":8000", nil)
}