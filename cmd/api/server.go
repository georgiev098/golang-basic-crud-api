package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "3000"

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the root route"))
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the teachers route"))
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the students route"))
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the execs  route"))
	})

	fmt.Println("Server running on port:", PORT)

	err := http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
