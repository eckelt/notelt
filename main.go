package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const FOLDER = "./notes"

func getNotePath(name string) string {
	path, _ := filepath.Abs(filepath.Join(FOLDER, filepath.Base(name+".md")))
	return path
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		// store request body in markdown file named after the note from the path in the notes folder

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read request body", http.StatusInternalServerError)
			return
		}
		noteName := r.PathValue("note")

		err = ioutil.WriteFile(getNotePath(noteName), body, 0644)

		// respond with 201 Created
		w.WriteHeader(http.StatusCreated)
	case "GET":
		noteName := r.PathValue("note")
		// read the file
		note, err := ioutil.ReadFile(getNotePath(noteName))
		if err != nil {
			http.Error(w, "could not read file", http.StatusInternalServerError)
			return
		}
		// write the file to the response
		w.Write(note)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	noteName := r.PathValue("note")
	// read the file
	note, err := ioutil.ReadFile(getNotePath(noteName))
	if err != nil {
		http.Error(w, "could not read file", http.StatusInternalServerError)
		return
	}
	// convert the markdown to html
	note = []byte(mdToHTML(note))
	// write the file to the response
	w.Write(note)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	//read index.html
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "could not read index.html", http.StatusInternalServerError)
		return
	}
	//write index.html to response
	w.Write(index)
}

func main() {
	fmt.Println("Starting server", getNotePath("test"))
	http.HandleFunc("/api/v1/note/{note}", apiHandler)
	http.HandleFunc("/{note}", markdownHandler)
	http.HandleFunc("/{note}/edit", editorHandler)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
