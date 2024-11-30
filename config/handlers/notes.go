package handlers

import (
	"fmt"
	"net/http"
	"your_project/models"
	"your_project/utils"
)

func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode note: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Insert the note into the database (ScyllaDB, Postgres, or any DB)
	err = utils.InsertNoteIntoDatabase(note)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert note: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Note created successfully!")
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := utils.GetNotesFromDatabase()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve notes: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}
