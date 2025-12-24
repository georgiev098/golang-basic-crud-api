package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/georgiev098/golang-basic-crud-api/internal/models"
	"github.com/georgiev098/golang-basic-crud-api/internal/repository/sqlconnect"
)

func AddTeacher(w http.ResponseWriter, r *http.Request) {

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "invalid request Body", http.StatusBadRequest)
		return
	}

	addedTeachers, err := sqlconnect.AddTeacherToDB(newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

func GetTeachers(w http.ResponseWriter, r *http.Request) {

	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDB(teachers, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetTeacher(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	// handle path param
	idNum, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid ID.", http.StatusBadRequest)
		return
	}

	teacher, err := sqlconnect.GetTeacherByIdDB(idNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher

	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacherDB(id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Conent-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}

func PatchTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]any

	err = json.NewDecoder(r.Body).Decode(&updates)

	updatedTeacher, err := sqlconnect.PatchSingleTeacherDB(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Conent-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)

}

func PatchTeachers(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload.", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchMultipleTeachersDB(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteSingleTeacherDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher deleted.",
		ID:     id,
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteTeachers(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteMultipleTeachersDB(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIds []int  `json:"deleted_ids"`
	}{
		Status:     "Teachers deleted.",
		DeletedIds: deletedIds,
	}

	json.NewEncoder(w).Encode(response)
}
