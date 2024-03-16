package main

import (
	"fmt"
	"io"
	"os"

	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type tmpls struct {
    tmpl *template.Template
}

func (t *tmpls) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.tmpl.ExecuteTemplate(w, name, data)
}

func init_tmpls() *tmpls {
    return &tmpls{
        tmpl: template.Must(template.ParseGlob("views/*.html")),
    }
}

//func give_root(w http.ResponseWriter, req *http.Request) {
//    tmpl := template.Must(template.ParseFiles("tmpl/index.html"))
//    tmpl.Execute(w, nil)
//}
//
//func add_item(w http.ResponseWriter, req *http.Request) {
//    data := req.PostFormValue("add-item-name")
//    tmpl := template.Must(template.ParseFiles("tmpl/item.html"))
//    tmpl.Execute(w, data)
////    sql_add_item(data)
//}
//
//func del_item(w http.ResponseWriter, req *http.Request) {
//    req.ParseForm()
////    data := req.PostFormValue("item-name")
//}
//

func root(c echo.Context) error {
    return c.Render(http.StatusOK, "index", nil)
}

var item_id = 0
type item struct {
    Name string
    Id int
}

func init_item(name string) item {
    item_id++
    return item {
        Name: name,
        Id: item_id,
    }
}

func add_item(c echo.Context) error {
    item := init_item(c.FormValue("add-item-name"))
    fmt.Printf("%+v\n", item)
    err := c.Render(http.StatusOK, "item", item)
    return err
}

func del_item(c echo.Context) error {
    return c.NoContent(http.StatusOK)
}

func sql_setup() {
    var err error
    db, err = sql.Open("mysql", "root:cumdump@(172.18.0.2:3306)/go_test?parseTime=true")
    if err != nil {
        fmt.Fprintf(os.Stderr, "setup_sql: %s\n", err)
    }
    err = db.Ping()
    if err != nil {
        fmt.Fprintf(os.Stderr, "setup_sql: Ping(): %s\n", err)
    }
}
//
//func sql_add_item(item string) {
//    _, err := db.Exec(`INSERT INTO items (name) VALUE (?)`, item)
//    if err != nil {
//        fmt.Fprintf(os.Stderr, "sql_add_item: %s\n", err)
//    }
//}

var db *sql.DB
func main() {
    e := echo.New()
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "method: ${method}, uri: ${uri}, status: ${status}\n",
    }))

    e.Renderer = init_tmpls()

    e.GET("/", root)
    e.POST("/items", add_item)
    e.DELETE("/items/:id", del_item)

    e.Static("/static", "static")

//    rt.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    sql_setup()

    e.Logger.Fatal(e.Start(":8080"))
}
