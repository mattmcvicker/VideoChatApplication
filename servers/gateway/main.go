package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
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

	topicsDirector := func(r *http.Request) {
		currentState := &handlers.SessionState{}
		_, err := sessions.GetState(r, context.SigningKey, context.SessionStore, currentState)
		if err != nil {
			fmt.Print("Not an authenticated user")
			r.Header.Del("X-User")
		} else {
			r.Header.Set("X-User", "{id: "+strconv.FormatInt(currentState.User.ID, 10)+"}")
		}

		serverName := "topics"
		r.Host = serverName
		r.URL.Host = serverName
		r.URL.Scheme = "https"
	}

	topicsProxy := &httputil.ReverseProxy{Director: topicsDirector}

	// define handlers for users and sessions resources
	mux.HandleFunc("/v1/users", context.UsersHandler)
	mux.HandleFunc("/v1/users/", context.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", context.SessionsHandler)
	mux.HandleFunc("/v1/sessions/", context.SpecificSessionHandler)
	mux.Handle("/v1/queue", topicsProxy)
	mux.Handle("/v1/topics", topicsProxy)

	// wrap api
	wrappedMux := handlers.NewCORS(mux)

	// start web server; report any errors
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}
