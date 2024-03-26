package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

func newContact(name, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}
func (d Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}
func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact("Brandon", "bb@gmail.com"),
			newContact("Kieluna", "kh@gmail.com"),
		},
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	page := newPage()
	e.Renderer = NewTemplates()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			FormData := newFormData()
			FormData.Values["name"] = name
			FormData.Values["email"] = email
			FormData.Errors["email"] = "Email already exists"
			return c.Render(422, "form", FormData)
		}

		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)
		c.Render(200, "form", newFormData())
		return c.Render(200, "oob-contact", contact)
	})
	e.Logger.Fatal(e.Start(":42069"))
}
