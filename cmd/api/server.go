package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/georgiev098/golang-basic-crud-api/internal/api/middleware"
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

	server := &http.Server{
		Addr:      ":" + PORT,
		Handler:   middleware.Compression(mux),
		TLSConfig: tlsConfig,
	}

	err := server.ListenAndServeTLS(cert, key)

	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

}
