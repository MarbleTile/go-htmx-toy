package main

import (
    "os"
	"fmt"
	"log"

	"net/http"
    "html/template"

    "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        log.Println(req.Method, req.URL.Path)
        f(w, req)
    }
}

func give_root(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodGet:
        tmpl := template.Must(template.ParseFiles("tmpl/index.html"))
        tmpl.Execute(w, nil)
        break
    default:
        return
    }
}

func add_item(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodPost:
        data := req.PostFormValue("add-item-name")
        tmpl := template.Must(template.ParseFiles("tmpl/item.html"))
        tmpl.Execute(w, data)
//        sql_add_item(data)
        break
    default:
        return
    }
}

func del_item(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case http.MethodPost:
        req.ParseForm()
//        data := req.PostFormValue("item-name")
        break
    default:
        return
    }
}

func sql_setup() {
    var err error
    db, err = sql.Open("mysql", "root:cumdump@(172.18.0.2:3306)/go_test?parseTime=true")
    if err != nil {
        fmt.Fprintf(os.Stderr, "setup_sql: %s", err)
    }
    err = db.Ping()
    if err != nil {
        fmt.Fprintf(os.Stderr, "setup_sql: Ping(): %s", err)
    }
}

func sql_add_item(item string) {
    _, err := db.Exec(`INSERT INTO items (name) VALUE (?)`, item)
    if err != nil {
        fmt.Fprintf(os.Stderr, "sql_add_item: %s", err)
    }
}

var db *sql.DB
func main() {
    http.HandleFunc("/",            logging(give_root))
    http.HandleFunc("/add_item",    logging(add_item))
    http.HandleFunc("/del_item",    logging(del_item))
    http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
        for k, v := range r.Header {
            fmt.Fprintf(w, "%s: %s\n", k, v)
        }
    })

    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static", http.StripPrefix("/static/", fs))

    sql_setup()

    http.ListenAndServe(":8080", nil)
}
