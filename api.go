package main

import (
	"goapi/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type APIServer struct {
	listenAddr string
	store      *PostgresStore
}

func NewApiServer(listenAddr string) *APIServer {
	store := NewPostgresStore()
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (a *APIServer) Run() {
	// start the scraper
	go func() {
		log.Printf("Scraper has started")
		startScraping(database.New(a.store.db), 10, time.Minute)
	}()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*, https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "UPDATE", "OPTIONS"},
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", a.handleReadiness)
	v1Router.Get("/error", a.handleError)
	v1Router.Post("/user", makeHTTPHandler(a.handleCreateUser))
	v1Router.Get("/user", makeHTTPHandler(a.middlewareAuth(a.handleGetUserByApiKey)))
	v1Router.Get("/users", makeHTTPHandler(a.handleGetUsers))
	v1Router.Post("/feed", makeHTTPHandler(a.middlewareAuth(a.handleCreateFeed)))
	v1Router.Get("/posts", makeHTTPHandler(a.middlewareAuth(a.handleGetPostsForUser)))
	v1Router.Get("/feeds", makeHTTPHandler(a.handleGetFeeds))
	v1Router.Post("/feedfollow", makeHTTPHandler(a.middlewareAuth(a.handleCreateFeedFollow)))
	v1Router.Get("/feedfollows", makeHTTPHandler(a.middlewareAuth(a.handleGetFeedFollows)))
	v1Router.Delete("/feedfollow/{feedID}", makeHTTPHandler(a.middlewareAuth(a.handleDeleteFeedFollow)))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + a.listenAddr,
		Handler: router,
	}

	log.Printf("Server is listening on PORT: %s", a.listenAddr)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server crashed!, Error: %s", err.Error())
	}
}
