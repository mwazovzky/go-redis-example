package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"go-redis-example/handlers"
)

type User struct {
	Name  string
	Email string
}

func main() {
	db := connectDB()
	defer db.Close()

	connectRedis()

	port := ":8080"
	router := mux.NewRouter()

	router.HandleFunc("/", Homepage)

	apiHandlers := handlers.NewApiHandlers(db)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", apiHandlers.Users).Methods(http.MethodGet)
	api.HandleFunc("/balances", apiHandlers.Balances).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         port,
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		log.Println("Starting http server at", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Println("Error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Recieved terminate signal, graceful shutdown, signal: [%s]", sig)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(ctx)

}

func Homepage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Welcome to Home Page")
}

func connectDB() *sql.DB {
	dbhost := "postgres"
	dbport := 5432
	dbuser := "user"
	dbpassword := "secret"
	dbname := "testdb"

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db")

	return db
}

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to redis")

	return client
}
