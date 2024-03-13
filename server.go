package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func give_root(w http.ResponseWriter, req *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/index.html"))
    tmpl.Execute(w, nil)
}

func add_item(w http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        return
    }
    data := req.PostFormValue("add-item-name")
    tmpl := template.Must(template.ParseFiles("templates/item.html"))
    tmpl.Execute(w, data)
}

func setup_sql() *sql.DB {
    db, err := sql.Open("mysql", "root:cumdump@(172.18.0.2:3306)/go_test?parseTime=true")
    if err != nil {
        fmt.Println("setup_sql:", err)
    }
    err = db.Ping()
    if err != nil {
        fmt.Println("db.Ping():", err)
    }
    return db
}

func main() {
    http.HandleFunc("/", give_root)
    http.HandleFunc("/add_item", add_item)

    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static", http.StripPrefix("/static/", fs))

    db := setup_sql()
    db.Ping()

    http.ListenAndServe(":8080", nil)
}
