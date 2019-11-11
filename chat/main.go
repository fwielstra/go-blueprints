package main

import (
	"flag"
	"go-blueprints/trace"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"

	// autoload .env file
	_ "github.com/joho/godotenv/autoload"
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
		t.templ = template.Must(template.ParseFiles(filepath.Join("chat/templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	// pass cookie value into template data for display.
	// TODO: XSS vulnerability unless template is secure. Still not very nice. Needs cookie signing and the like to detect tampering.
	// NOTE: No-op if no auth cookie is set; I guess it's not an error state to not be authenticated, so fine to ignore
	// the err.
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	flag.Parse()

	gomniauth.SetSecurityKey(os.Getenv("SESSION_SECRET"))
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_OAUTH_KEY"), os.Getenv("GOOGLE_OAUTH_SECRET"), "http://localhost:8080/auth/callback/google"),
	)

	// create room
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// set up HTTP routes
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
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
