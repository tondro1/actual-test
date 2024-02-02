package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/cors"
	"github.com/tondro1/actual-test/internal/database"
)

type apiCfg struct {
	db *database.Queries
}
const DB_URL = "postgres://alexchoi:@localhost:5432/main-go?sslmode=disable"
func main() {
	// compile templates
	// templates := map[string]*template.Template{}
	// tmpIndex := template.Must(template.ParseFiles("./public/root.html"))
	// templates["index"] = tmpIndex
	conn, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(context.Background())

	db := apiCfg{db: database.New(conn)}

	pageRouter := chi.NewRouter()
	apiRouter := chi.NewRouter()
	
	// Data API
	apiRouter.Use(cors.AllowAll().Handler)
	apiRouter.Post("/register", register)

	go func() {
		log.Println("Starting server on localhost:1324")
		http.ListenAndServe(":1324", apiRouter)
	}()

	// HTML API

	// static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	jsFs := http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js")))

	pageRouter.Handle("/static/*", fs)
	pageRouter.HandleFunc("/favicon.ico", handlerFavicon)
	pageRouter.Handle("/js/*", jsFs)
	
	pageRouter.Get("/", renderIndex)
	pageRouter.Get("/login", renderLogin)
	pageRouter.Get("/register", renderRegister)

	log.Println("Starting server on localhost:1323")
	http.ListenAndServe(":1323", pageRouter)
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}