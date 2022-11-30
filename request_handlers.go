package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Comment `json:"data"`
	Message string    `json:"message"`
}

// Retrieve all the comments
func getComments(w http.ResponseWriter, r *http.Request) {
	var comments = Comment{}
	var comments_slice = []Comment{}

	db := dbConn()

	row, err := db.Query("SELECT * FROM comments")
	if err != nil {
		println("Table not found.")
	}
	for row.Next() {
		err := row.Scan(&comments.Id, &comments.Title, &comments.Text, &comments.Author, &comments.Date, &comments.Anime)
		if err != nil {
			println("Couldn't find searched params.")
			log.Fatal(err)
		}

		comments_slice = append(comments_slice, comments)
	}

	var response = JsonResponse{Type: "success", Data: comments_slice}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}

// Retrieve single comment by ID
func getCommentByID(w http.ResponseWriter, r *http.Request) {
	var comments = Comment{}

	// Get the ID passed as a param in the url parsing the URL
	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/comment/"), 10, 64)
	if err != nil {
		println("Error parsing through the url.")
		log.Fatal(err)
	}

	db := dbConn()
	row := db.QueryRow("SELECT * FROM comments WHERE id = ?", id)
	if err := row.Scan(&comments.Id, &comments.Title, &comments.Text, &comments.Author, &comments.Date, &comments.Anime); err != nil {
		if err == sql.ErrNoRows {
			println("Row no found.")
			log.Fatal(err)
		}
		println("Couldn't find searched params.")
		log.Fatal(err)
	}

	var response = JsonResponse{Type: "success", Data: []Comment{comments}}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}

// Post new comment
func postComment(w http.ResponseWriter, r *http.Request) {
	var newComment Comment

	db := dbConn()
	decoder := json.NewDecoder(r.Body).Decode(&newComment)
	if err := decoder; err != nil {
		println("Couldn't decode comment.")
		log.Fatal(err)
	}

	row, err := db.Prepare("INSERT INTO comments(title, comment_text, author, publish_date, anime) VALUES(?,?,?,?,?)")
	if err != nil {
		println("Error preparing the mysql command.")
		log.Fatal(err)
	}

	res, err := row.Exec(newComment.Title, newComment.Text, newComment.Author, newComment.Date, newComment.Anime)
	if err != nil {
		println("Error executing the mysql insertion.")
		log.Fatal(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		println("Error finding the id of the last insertion to the database.")
		log.Fatal(err)
	}

	newComment.Id = lastId

	// Send back a status "created" with the newComment json file
	var response = JsonResponse{Type: "success", Data: []Comment{newComment}}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}

// Delete comment
func deleteComment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/comment/delete/")

	db := dbConn()
	delete, err := db.Prepare("DELETE FROM comments WHERE id=?")
	if err != nil {
		println("Error preparing the mysql delete.")
		log.Fatal(err)
	}
	delete.Exec(id)

	var response = JsonResponse{Type: "success"}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}

// Update comment
func updateComment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/comment/update/")
	var editedComment Comment

	decoder := json.NewDecoder(r.Body).Decode(&editedComment)
	if err := decoder; err != nil {
		return
	}

	db := dbConn()
	if editedComment.Title != "" {
		edit, err := db.Prepare("UPDATE comments SET title=? WHERE id=?")
		if err != nil {
			println("Error preparing the mysql update.")
			log.Fatal(err)
		}
		edit.Exec(editedComment.Title, id)
	}
	if editedComment.Text != "" {
		edit, err := db.Prepare("UPDATE comments SET comment_text=? WHERE id=?")
		if err != nil {
			println("Error preparing the mysql update.")
			log.Fatal(err)
		}
		edit.Exec(editedComment.Text, id)
	}

	var response = JsonResponse{Type: "success"}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}
