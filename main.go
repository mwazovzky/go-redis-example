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

	"go-redis-example/http/handlers"
	"go-redis-example/http/middleware"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var port string

func init() {
	godotenv.Load()
	port = fmt.Sprintf(":%s", os.Getenv("PORT"))
}

func main() {
	db := connectDB()
	defer db.Close()

	cache := connectRedis()

	port := ":8080"
	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/", Homepage)

	apiHandlers := handlers.NewApiHandlers(db, cache)
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/users", apiHandlers.Users).Methods(http.MethodGet)
	api.HandleFunc("/balances", apiHandlers.Balances).Methods(http.MethodGet)
	api.HandleFunc("/clear", apiHandlers.Clear).Methods(http.MethodGet)

	router.Use(middleware.NewMetricsMiddleware().Metrics)

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
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	db, err := sql.Open("postgres", dsn)
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
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
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
