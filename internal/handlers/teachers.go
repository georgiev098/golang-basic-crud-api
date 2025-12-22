package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/georgiev098/golang-basic-crud-api/internal/models"
	"github.com/georgiev098/golang-basic-crud-api/internal/repository/sqlconnect"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

// func init() {
// 	teachers[nextId] = models.Teacher{
// 		ID:        nextId,
// 		FirstName: "John",
// 		LastName:  "Doe",
// 		Class:     "9A",
// 		Subject:   "Math",
// 	}
// 	nextId++
// 	teachers[nextId] = models.Teacher{
// 		ID:        nextId,
// 		FirstName: "Jane",
// 		LastName:  "Doe",
// 		Class:     "10A",
// 		Subject:   "Algebra",
// 	}

// }

func addTeacher(w http.ResponseWriter, r *http.Request) {
	// connect to DB
	db, err := sqlconnect.ConnectToDB("school")
	if err != nil {
		http.Error(w, "Could not establish DB connection.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request Body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Error in preparing DB query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		resp, err := stmt.Exec(
			newTeacher.FirstName,
			newTeacher.LastName,
			newTeacher.Email,
			newTeacher.Class,
			newTeacher.Subject,
		)
		if err != nil {
			http.Error(w, "Error inserting data into DB.", http.StatusInternalServerError)
			return
		}

		newId, err := resp.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting newly created ID.", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(newId)
		addedTeachers[i] = newTeacher

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `jsom:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(resp)

}

func getTeacher(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectToDB("school")
	if err != nil {
		http.Error(w, "Could not establish DB connection.", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []interface{}

		if firstName != "" {
			query += " AND first_name = ?"
			args = append(args, firstName)
		}

		if lastName != "" {
			query += " AND last_name = ?"
			args = append(args, lastName)
		}

		rows, err := db.Query(query, args...)

		if err != nil {
			http.Error(w, "Database query error.", http.StatusInternalServerError)
			fmt.Print(err)
		}

		defer rows.Close()

		teacherList := make([]models.Teacher, 0)
		for rows.Next() {
			var teacher models.Teacher

			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(w, "Database scanning db results.", http.StatusInternalServerError)
				fmt.Print(err)
			}

			teacherList = append(teacherList, teacher)
		}

		resp := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
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

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", idNum).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeacher(w, r)
	case http.MethodPost:
		addTeacher(w, r)
	}
}
