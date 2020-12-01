package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mattmcvicker/WeFeud/servers/gateway/handlers"
	"github.com/mattmcvicker/WeFeud/servers/gateway/models/users"
	"github.com/mattmcvicker/WeFeud/servers/gateway/sessions"
)

//main is the main entry point for the server
func main() {
	// read in the ADDR environment variable
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}
	// read in the session key
	sessKey := os.Getenv("SESSIONKEY")
	// read in the redisaddr
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}
	// read in full data source name
	dsn := os.Getenv("DSN")

	// create new mux for the web server
	mux := http.NewServeMux()

	// create db and handler context
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	redisStore := sessions.NewRedisStore(client, time.Hour)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	sqlStore := users.NewSQLStore(db)

	context := handlers.NewHandlerContext(sessKey, redisStore, sqlStore)

	mux2 := http.NewServeMux()

	// define handlers for users and sessions resources
	mux2.HandleFunc("/v1/users", context.UsersHandler)
	mux2.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux2.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux2.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)

	// wrap api
	wrappedMux := handlers.NewCORS(mux2)

	// add to main mux
	mux.Handle("/v1/users", wrappedMux)
	mux.Handle("/v1/users/", wrappedMux)
	mux.Handle("/v1/sessions", wrappedMux)
	mux.Handle("/v1/sessions/", wrappedMux)

	// start web server; report any errors
	log.Fatal(http.ListenAndServe(addr, mux))
}
