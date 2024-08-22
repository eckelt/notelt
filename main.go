package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const notes_folder = "./notes"

func getNotePath(name string) string {
	path, _ := filepath.Abs(filepath.Join(notes_folder, filepath.Base(name+".md")))
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
			http.Error(w, "could not read file node:"+getNotePath(noteName), http.StatusInternalServerError)
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
		http.Error(w, "could not read file"+getNotePath(noteName), http.StatusInternalServerError)
		return
	}
	// convert the markdown to html
	note = []byte(mdToHTML(note))
	template, err := ioutil.ReadFile("static/template.html")
	if err != nil {
		http.Error(w, "could not read template.html", http.StatusInternalServerError)
		return
	}
	template = []byte(strings.Replace(string(template), "<!-- page content -->", string(note), 1))
	// write the file to the response
	w.Write(template)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	//read edit.html
	edit, err := ioutil.ReadFile("edit.html")
	if err != nil {
		http.Error(w, "could not read edit.html", http.StatusInternalServerError)
		return
	}
	//write edit.html to response
	w.Write(edit)
}

func main() {
	// if os.Getenv("NOTES_FOLDER") == "" {
	// 	log.Fatal("NOTES_FOLDER environment variable must be set")
	// }
	// notes_folder = os.Getenv("NOTES_FOLDER")
	fmt.Println("Starting server", getNotePath(""))
	http.HandleFunc("/api/v1/note/{note}", apiHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/{note}", markdownHandler)
	http.HandleFunc("/e/{note}", editorHandler)
	fmt.Println("Listening on port 8520")
	http.ListenAndServe(":8520", nil)
}
