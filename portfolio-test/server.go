package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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


type model_data struct {
    Name string
}

type model_slc struct {
    Models []model_data
}

func (m *model_slc) add_model(model model_data){
    m.Models = append(m.Models, model)
}

var model_list model_slc
func init_model_list() {
    fs.WalkDir(os.DirFS("static/models"), ".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            fmt.Println(err)
        }
        if path != "." {
            model_list.add_model(model_data{Name: path[:strings.IndexByte(path, '.')]})
        }
        return nil
    })
}

func init_smtp() {
    
}

func give_root(c echo.Context) error {
    return c.Render(http.StatusOK, "index", nil)
}

func give_model_list(c echo.Context) error {
    b, err := json.Marshal(model_list.Models)
    if err != nil {
        return c.String(http.StatusNotFound, "no models")
    }
    return c.JSONBlob(http.StatusOK, b)
}

func give_about(c echo.Context) error {
    return c.Render(http.StatusOK, "about", nil)
}

func give_contact(c echo.Context) error {
    return c.Render(http.StatusOK, "contact", nil)
}

func give_work(c echo.Context) error {
    return c.String(http.StatusOK, "")
}

func main() {
    init_model_list()

    e := echo.New()
    e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
        Format: "method: ${method}, uri: ${uri}, status: ${status}\n",
    }))

    e.Renderer = init_tmpls()

    e.GET("/", give_root)
    e.GET("/model_list", give_model_list);
    e.GET("/about", give_about);
    e.GET("/contact", give_contact);
    e.GET("/work", give_work);
    e.Static("/static", "static")

    e.Logger.Fatal(e.Start(":8080"))
}
