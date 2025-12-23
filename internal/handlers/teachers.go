package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/georgiev098/golang-basic-crud-api/internal/models"
	"github.com/georgiev098/golang-basic-crud-api/internal/repository/sqlconnect"
)

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

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	return validFields[field]

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
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"

		var args []any

		query, args = addFilters(r, query, args)

		query = addSorting(r, query)

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

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sort-by"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]

			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}

			query += " " + field + " " + order
		}
	}
	return query
}

func updateTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher

	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)

	db, err := sqlconnect.ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, subject, class FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.Class, &existingTeacher.Email, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Teacher not found.", http.StatusNotFound)
		return
	} else {
		if err != nil {
			log.Println(err)
			http.Error(w, "Retrieving teacher from DB.", http.StatusInternalServerError)
			return
		}
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating entry.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Conent-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

func patchTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)

	db, err := sqlconnect.ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, subject, class FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.Class, &existingTeacher.Email, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Teacher not found.", http.StatusNotFound)
		return
	} else {
		if err != nil {
			log.Println(err)
			http.Error(w, "Retrieving teacher from DB.", http.StatusInternalServerError)
			return
		}
	}

	teacherVal := reflect.ValueOf(existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			field.Tag.Get("json")

			if field.Tag.Get("json") == k+" ,omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}

		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating entry.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Conent-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)

}

func addFilters(r *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}

	return query, args
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeacher(w, r)
	case http.MethodPost:
		addTeacher(w, r)
	case http.MethodPut:
		updateTeacher(w, r)
	}
}
