package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
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
	// read in topic and quiz microservice address
	feudAddress := os.Getenv("FEUDADDR")
	// read in the session key
	sessKey := os.Getenv("SESSIONKEY")
	// read in the redisaddr
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
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

	// create new mux for cors middleware
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
	//mux.HandleFunc("/v1/queue", )
	mux.HandleFunc("/test", handlers.HandleTestPath)

	// define reverse proxy for topic and quiz microservice
	feudDirector := func(r *http.Request) {
		// check for authenticated user and add X-User header (with id)
		if r.Header.Get("X-User") != "" {
			r.Header.Set("X-User", "")
		}
		user, _ := handlers.GetAuthenticatedUser(context, r)
		if user != nil {
			// encode user into json
			uJSON, _ := json.Marshal(user)
			r.Header.Add("X-User", string(uJSON))
		}

		r.Host = feudAddress
		r.URL.Host = feudAddress
		r.URL.Scheme = "http"
	}
	feudProxy := &httputil.ReverseProxy{Director: feudDirector}
	mux.Handle("/v1/topics", feudProxy)
	mux.Handle("/v1/topics/", feudProxy)
	mux.Handle("/v1/queue", feudProxy)

	// start web server; report any errors
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
