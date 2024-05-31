package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

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

func viewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Remove the /view/ from the request url and keep everything until the end of the string (using :)
	title := request.URL.Path[len("/view/"):]
	page, _ := loadPage(title)

	renderTemplate(responseWriter, "view", page)
}

func editHandler(responseWriter http.ResponseWriter, request *http.Request) {
	title := request.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	renderTemplate(responseWriter, "edit", page)
}

func saveHandler(responseWriter http.ResponseWriter, request *http.Request) {
}

func renderTemplate(responseWritter http.ResponseWriter, fileName string, page *Page) {
	_template, _ := template.ParseFiles(fileName + ".html")
	_template.Execute(responseWritter, page)
}
