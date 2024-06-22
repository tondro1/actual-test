package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/tondro1/actual-test/api"
	"github.com/tondro1/actual-test/internal/database"
)

func main() {
	// compile templates
	// templates := map[string]*template.Template{}
	// tmpIndex := template.Must(template.ParseFiles("./public/root.html"))
	// templates["index"] = tmpIndex

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	DB_URL := os.Getenv("DB_URL")
	PORT := os.Getenv("PORT")

	conn, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(context.Background())

	apiCfg := api.ApiCfg{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	
	// router.Use(cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:1323"},
	// 	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders: []string{"*"},
	// 	ExposedHeaders: []string{"Link"},
	// 	MaxAge: 300,
	// 	AllowCredentials: true,
	// }).Handler)
	router.Post("/api/register", apiCfg.Register)
	router.Post("/api/login", apiCfg.Login)
	router.Put("/api/logout", apiCfg.Validate(apiCfg.Authenticate(apiCfg.Logout)))


	// static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	jsFs := http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js")))

	router.Handle("/static/*", fs)
	router.HandleFunc("/favicon.ico", handlerFavicon)
	router.Handle("/js/*", jsFs)
	
	router.Get("/", apiCfg.Validate(apiCfg.Authenticate(api.RenderIndex)))
	router.Get("/login", apiCfg.Validate(apiCfg.Authenticate(api.RenderLogin)))
	router.Get("/register", apiCfg.Validate(apiCfg.Authenticate(api.RenderRegister)))

	log.Println("Starting server on localhost:" + PORT)
	http.ListenAndServe(":" + PORT, router)
}

func handlerFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}