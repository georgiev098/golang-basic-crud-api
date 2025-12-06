package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/georgiev098/golang-basic-crud-api/internal/api/middleware"
	"github.com/georgiev098/golang-basic-crud-api/internal/middlewares"
	"github.com/georgiev098/golang-basic-crud-api/internal/repository/sqlconnect"
	"github.com/georgiev098/golang-basic-crud-api/internal/router"
	"github.com/georgiev098/golang-basic-crud-api/pkg/utils"
)

const PORT = "3000"

func main() {
	// connect to DB
	db, err := sqlconnect.ConnectToDB("school")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cert := "certs/localhost.crt"
	key := "certs/localhost.key"

	router := router.Rotuer()

	fmt.Println("Server running on port:", PORT)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rl := middlewares.NewRateLimiter(5, time.Minute)

	hppOptions := middlewares.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-from-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	secureMux := utils.ApplyMiddlewares(router, middlewares.Hpp(hppOptions), middleware.Compression)

	server := &http.Server{
		Addr:      ":" + PORT,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS(cert, key)

	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

}
