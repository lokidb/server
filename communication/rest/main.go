package rest

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lokidb/engine"
)

type server struct {
	host   string
	port   int
	addr   string
	router *http.Handler
}

func loadRouters(r *mux.Router, engine *engine.KeyValueStore) {
	r.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		handleGet(w, r, engine)
	}).Name("get")

	r.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		handleSet(w, r, engine)
	}).Name("set")

	r.HandleFunc("/del", func(w http.ResponseWriter, r *http.Request) {
		handleDel(w, r, engine)
	}).Name("del")

	r.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		handleKeys(w, r, engine)
	}).Name("keys")

	r.HandleFunc("/flush", func(w http.ResponseWriter, r *http.Request) {
		handleFlush(w, r, engine)
	}).Name("flush")
}

func NewServer(host string, port int, engine *engine.KeyValueStore) *server {
	s := new(server)
	s.host = host
	s.port = port
	s.addr = s.host + ":" + strconv.Itoa(s.port)

	r := mux.NewRouter()
	loadRouters(r, engine)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	s.router = &loggedRouter

	return s
}

func (s server) Run() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	srv := &http.Server{
		Addr: s.addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      *s.router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("REST server listening at %s\n", s.addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
