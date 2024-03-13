package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"
)

func headers(w http.ResponseWriter, req *http.Request){
    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func give_root(w http.ResponseWriter, req *http.Request){
    tmpl := template.Must(template.ParseFiles("templates/index.html"))
    tmpl.Execute(w, nil)
}

type item struct {
    name string
}
func add_item(w http.ResponseWriter, req *http.Request){
    if req.Method != "POST" {
        return
    }
    data := req.PostFormValue("add-item-name")
    fmt.Println(data)
    tmpl := template.Must(template.ParseFiles("templates/item.html"))
    time.Sleep(1 * time.Second)
    tmpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/headers", headers)

    http.HandleFunc("/", give_root)
    http.HandleFunc("/add_item", add_item)

    http.ListenAndServe(":8080", nil)
}
