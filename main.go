package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/CODECZERO/goalngtest/internal/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	godotenv.Load() //takes env from file
	//get's port number from env
	port := os.Getenv("PORT")
	if port == " " { //checks if it's empty if it's then log error
		log.Fatal("port number is missing")
	}
	//get's database url env
	dburl := os.Getenv("DBURL")
	if dburl == " " { //checks if it's empty if it's then log error
		log.Fatal("port number is missing")
	}

	conn, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Fatal("something went wrong while connection to database", err)
	}

	dbQueris := db.New(conn)
	if err != nil {
		log.Fatal("something went wrong while runing query in database", err)
	}

	apiCfg := apiConfig{
		DB: dbQueris,
	}

	router := chi.NewRouter() //create router for go lang using chi router

	router.Use(cors.Handler(cors.Options{ //config of cors
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()     //main router or backend url router
	v1router.Get("/healt", handler) //allows get methode to check wheater app is running or not
	v1router.Post("/user", apiCfg.HandlerUser)
	v1router.Get("/user", apiCfg.middleware(apiCfg.HandlerGetUser))
	v1router.Post("/feed", apiCfg.middleware(apiCfg.handlerCreateFeed))
	v1router.Get("/feed", apiCfg.handlerGetFeeds)

	v1router.Get("/feedFollower",apiCfg.middleware((apiCfg.handlerFeedFollowsGet)))
	v1router.Post("feedFollower",apiCfg.middleware(apiCfg.handlerCreateFeed))
	v1router.Delete("/feedFollower",apiCfg.middleware(apiCfg.handlerDeleteFollwer))
	

	router.Mount("/v1", v1router) //is main url of /healt or /v1/healt

	server := &http.Server{ //making server
		Handler: router, //passing router to main server
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe()) //starting server at given port
}
