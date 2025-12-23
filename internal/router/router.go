package router

import (
	"net/http"

	"github.com/georgiev098/golang-basic-crud-api/internal/handlers"
)

func Rotuer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("GET /teachers/", handlers.GetTeachers)
	mux.HandleFunc("POST /teachers/", handlers.AddTeacher)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacher)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacher)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacher)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacher)

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
