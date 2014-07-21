package web

import (
	"encoding/xml"
	"github.com/unrolled/render"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Renderer interface {
	XML(w http.ResponseWriter, status int, v interface{})
	HTML(w http.ResponseWriter, status int, name string, binding interface{})
}

type renderer struct {
	renderer *render.Render
}

func (r *renderer) HTML(w http.ResponseWriter, status int, name string, binding interface{}) {
	r.renderer.HTML(w, status, name, binding)
}

func (r *renderer) XML(w http.ResponseWriter, status int, v interface{}) {
	r.renderer.XML(w, status, v)
}

func getDateString(d time.Time) string {
	if d.IsZero() {
		return "Date TBC"
	} else {
		year, month, day := d.Date()
		suffix := "th"
		switch day % 10 {
		case 1:
			if day%100 != 11 {
				suffix = "st"
			}
		case 2:
			if day%100 != 12 {
				suffix = "nd"
			}
		case 3:
			if day%100 != 13 {
				suffix = "rd"
			}
		}
		return d.Weekday().String() + " " + strconv.Itoa(day) + suffix + " " + month.String() + " " + strconv.Itoa(year)
	}
}

func NewRenderer() Renderer {
	r := render.New(render.Options{
		Layout:    "index",
		Delims:    render.Delims{"[[", "]]"},
		PrefixXML: []byte(xml.Header),
		Funcs: []template.FuncMap{
			{
				"AsHTML": func(s string) template.HTML {
					return template.HTML(s)
				},
				"AsDate": getDateString,
			},
		},
	})
	return &renderer{renderer: r}
}
