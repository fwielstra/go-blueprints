package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// commandline flags
var (
	addr = flag.String("addr", ":8080", "The addr of the application")
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	flag.Parse()

	// create room
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)

	// set up HTTP routes
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	// start the room
	go r.run()

	// start the web server
	log.Println("Starting chat server on", *addr)
	// actually does ListenAndServe continue running until err?
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
