package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/georgiev098/golang-basic-crud-api/internal/api/middleware"
	"github.com/georgiev098/golang-basic-crud-api/internal/middlewares"
)

const PORT = "3000"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the root route"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the teachers route"))
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the students route"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the execs  route"))
}

func main() {

	cert := "certs/localhost.crt"
	key := "certs/localhost.key"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server running on port:", PORT)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := middlewares.NewRateLimiter(5, time.Minute)

	hppOptions := middlewares.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-from-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	secureMux := middlewares.Hpp(hppOptions)(rl.Middleware(middleware.Compression(mux)))

	server := &http.Server{
		Addr:      ":" + PORT,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	err := server.ListenAndServeTLS(cert, key)

	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

}
