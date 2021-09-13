package main

import (
	"chat/auth"
	"chat/trace"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.template.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	log.Println("Starting web server on", *addr)
	mux := http.NewServeMux()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	mux.Handle("/", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	mux.Handle("/chat/{room_id}", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	mux.Handle("/room/1", r)
	mux.Handle("/login", &templateHandler{filename: "login.html"})
	go r.run()
	if err := http.ListenAndServe(*addr, auth.New(mux)); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
