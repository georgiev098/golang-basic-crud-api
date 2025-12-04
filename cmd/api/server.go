package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/georgiev098/golang-basic-crud-api/internal/api/middleware"
	"github.com/georgiev098/golang-basic-crud-api/internal/middlewares"
)

const PORT = "3000"

type Teacher struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"sbuject,omitempty"`
}

var (
	teachers = make(map[int]Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Algebra",
	}

}

func getTeacher(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		teacherList := make([]Teacher, 0, len(teachers))

		for _, v := range teachers {
			if (firstName == "" || v.FirstName == firstName) && (lastName == "" || v.LastName == lastName) {
				teacherList = append(teacherList, v)
			}
		}

		resp := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}

	// handle path param
	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	teacher, exists := teachers[idNum]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(teacher)
}

func addTeacher(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request Body", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextId
		teachers[nextId] = newTeacher
		addedTeachers[i] = newTeacher
		nextId++
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `jsom:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(resp)

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the root route"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeacher(w, r)
	case http.MethodPost:
		addTeacher(w, r)
	}
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

	// rl := middlewares.NewRateLimiter(5, time.Minute)

	hppOptions := middlewares.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-from-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	secureMux := applyMiddlewares(mux, middlewares.Hpp(hppOptions), middleware.Compression)

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

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
