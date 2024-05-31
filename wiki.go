package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

var (
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

// Method added to the Page type
func (p *Page) save() error {
	filename := p.Title + ".txt"
	// 0600 - Create the file with read-write permissions for the current user
	return os.WriteFile(filename, p.Body, 0600)
}

// This is a function as it's not attached to a type
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func viewHandler(responseWriter http.ResponseWriter, request *http.Request, title string) {
	// Remove the /view/ from the request url and keep everything until the end of the string (using :)
	// title := request.URL.Path[len("/view/"):]
	title, err := getTitle(responseWriter, request)
	if err != nil {
		return
	}

	page, err := loadPage(title)
	if err != nil {
		http.Redirect(responseWriter, request, "/edit"+title, http.StatusFound)
		return
	}

	renderTemplate(responseWriter, "view", page)
}

func editHandler(responseWriter http.ResponseWriter, request *http.Request, title string) {
	title, err := getTitle(responseWriter, request)
	if err != nil {
		return
	}

	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(responseWriter, "edit", page)
}

func saveHandler(responseWriter http.ResponseWriter, request *http.Request, title string) {
	title, err := getTitle(responseWriter, request)
	if err != nil {
		return
	}

	body := request.FormValue("body")

	page := &Page{Title: title, Body: []byte(body)}

	err = page.save()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(responseWriter, request, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	// return func(responseWriter http.responseWriter, request *http.Request) {
	// }
}

func getTitle(responseWriter http.ResponseWriter, request *http.Request) (string, error) {
	match := validPath.FindStringSubmatch(request.URL.Path)

	if match == nil {
		http.NotFound(responseWriter, request)
		return "", errors.New("invalid page title")
	}
	return match[2], nil
}

func renderTemplate(responseWritter http.ResponseWriter, fileName string, page *Page) {
	err := templates.ExecuteTemplate(responseWritter, fileName+".html", page)
	if err != nil {
		http.Error(responseWritter, err.Error(), http.StatusInternalServerError)
	}
}
