package main

import (
	"fmt"
	"log"
	"net/http"
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

func main() 

	cert := "cert.pem"
	key := "key.pem"

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server running on port:", PORT)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("Error starting the server", err)
	}

