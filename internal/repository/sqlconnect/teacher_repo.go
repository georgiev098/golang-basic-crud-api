package sqlconnect

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/georgiev098/golang-basic-crud-api/internal/models"
)

func GetTeachersDB(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectToDB("school")
	if err != nil {
		// http.Error(w, "Could not establish DB connection.", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"

	var args []any

	query, args = AddFilters(r, query, args)

	query = AddSorting(r, query)

	rows, err := db.Query(query, args...)

	if err != nil {
		// http.Error(w, "Database query error.", http.StatusInternalServerError)
		fmt.Print(err)
	}

	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher

		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

		if err != nil {
			// http.Error(w, "Database scanning db results.", http.StatusInternalServerError)
			fmt.Print(err)
			return nil, err
		}

		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func GetTeacherByIdDB(idNum int) (models.Teacher, error) {
	db, err := ConnectToDB("school")
	if err != nil {
		// http.Error(w, "Could not establish DB connection.", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", idNum).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	} else if err != nil {
		// http.Error(w, "Database query error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return teacher, nil
}

func AddTeacherToDB(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectToDB("school")
	if err != nil {
		// http.Error(w, "Could not establish DB connection.", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		// http.Error(w, "Error in preparing DB query", http.StatusInternalServerError)
		return nil, err
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
			// http.Error(w, "Error inserting data into DB.", http.StatusInternalServerError)
			return nil, err
		}

		newId, err := resp.LastInsertId()
		if err != nil {
			// http.Error(w, "Error getting newly created ID.", http.StatusInternalServerError)
			return nil, err
		}
		newTeacher.ID = int(newId)
		addedTeachers[i] = newTeacher

	}
	return addedTeachers, nil
}

func UpdateTeacherDB(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, subject, class FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.Class, &existingTeacher.Email, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println(err)
		// http.Error(w, "Teacher not found.", http.StatusNotFound)
		return models.Teacher{}, err
	} else {
		if err != nil {
			log.Println(err)
			// http.Error(w, "Retrieving teacher from DB.", http.StatusInternalServerError)
			return models.Teacher{}, err
		}
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error updating entry.", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return updatedTeacher, nil
}

func PatchMultipleTeachersDB(updates []map[string]any) error {
	db, err := ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error starting transaction.", http.StatusInternalServerError)
		return err
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			// http.Error(w, "Invalid integer.", http.StatusBadRequest)
		}

		id, err := strconv.Atoi(idStr)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// http.Error(w, "Error converting ID to Int.", http.StatusNotFound)
				return err
			}
		}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ? ", id).Scan(&teacherFromDb.ID, &teacherFromDb.Class, &teacherFromDb.Email, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// http.Error(w, "Teacher not found.", http.StatusNotFound)
				return err
			} else {
				log.Println(err)
				// http.Error(w, "Could not query row.", http.StatusInternalServerError)
				return err
			}
		}

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)

				if field.Tag.Get("json") == k+" ,omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v", val.Type(), fieldVal.Type())
							return err
						}
					}
					break
				}

			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, class = ?, email = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Class, teacherFromDb.Email, teacherFromDb.Subject)
		if err != nil {
			// http.Error(w, "Error updating teacher.", http.StatusInternalServerError)
			tx.Rollback()
		}
	}

	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Could not commit changes.", http.StatusInternalServerError)
		tx.Rollback()
	}
	return nil
}

func PatchSingleTeacherDB(id int, updates map[string]any) (models.Teacher, error) {
	db, err := ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, subject, class FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.Class, &existingTeacher.Email, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		log.Println(err)
		// http.Error(w, "Teacher not found.", http.StatusNotFound)
		return models.Teacher{}, err
	} else {
		if err != nil {
			log.Println(err)
			// http.Error(w, "Retrieving teacher from DB.", http.StatusInternalServerError)
			return models.Teacher{}, err
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
		// http.Error(w, "Error updating entry.", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return existingTeacher, nil
}

func DeleteSingleTeacherDB(id int) error {
	db, err := ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return err
	}

	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ? ", id)
	if err != nil {
		log.Println(err)
		// http.Error(w, "Could not delete teacher.", http.StatusInternalServerError)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error retrieving delete result.", http.StatusInternalServerError)
		return err
	}

	if rowsAffected == 0 {
		// http.Error(w, "Teacher not found.", http.StatusNotFound)
		return err
	}
	return nil
}

func DeleteMultipleTeachersDB(ids []int) ([]int, error) {
	db, err := ConnectToDB("teachers")
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error starting transaction to DB", http.StatusInternalServerError)
		return nil, err
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		log.Println(err)
		// http.Error(w, "Could not prepare DB statement.", http.StatusInternalServerError)
		return nil, err
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		result, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			// http.Error(w, "Error deleting teacher.", http.StatusInternalServerError)
			return nil, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			log.Println(err)
			// http.Error(w, "Error retrieving deleted results.", http.StatusInternalServerError)
			return nil, err
		}

		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

		if rowsAffected < 1 {
			tx.Rollback()
			// http.Error(w, fmt.Sprintf("ID %d does not exists", id), http.StatusInternalServerError)
			return nil, err
		}

	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		// http.Error(w, "Error commiting results.", http.StatusInternalServerError)
		return nil, err
	}

	if len(deletedIds) < 1 {
		// http.Error(w, "Ids do not exist.", http.StatusInternalServerError)
		return nil, err
	}
	return deletedIds, nil
}

func AddFilters(r *http.Request, query string, args []any) (string, []any) {
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

func AddSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sort-by"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field, order := parts[0], parts[1]

			if !IsValidSortField(field) || !IsValidSortOrder(order) {
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

func IsValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func IsValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	return validFields[field]

}
