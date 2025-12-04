package handlers

import "net/http"

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the students route"))
}
