package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Comments data structure.

type Comment struct {
	ID     int64
	Title  string
	Text   string
	Author string
	Date   string
	Anime  string
}

// Json file that will be send

type JsonResponseComments struct {
	Type    string    `json:"type"`
	Data    []Comment `json:"data"`
	Message string    `json:"message"`
}

// Main function

func main() {
	// Endpoint definitions
	// Associate the endpoint "/comments" with the getComments function.
	http.HandleFunc("/comments", getComments)

	// Associate the endpoint "/comments" with the getCommentByID function.
	http.HandleFunc("/comment/", getCommentByID)

	// Associate the endpoint "/comments" with the postComment function.
	http.HandleFunc("/new_comment", postComment)

	// Associate the endpoint "/comments" with the deleteComment function.
	http.HandleFunc("/comment/delete/", deleteComment)

	// Associate the endpoint "/comments" with the deleteComment function.
	http.HandleFunc("/comment/update/", updateComment)

	// Start the associate the router whit one http server
	fmt.Println("Server at 8080")
	http.ListenAndServe(":8080", nil)
}

// Seeding comments data

var comments = []Comment{
	{ID: 1, Title: "Testing Forum", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Darling"},
	{ID: 2, Title: "Testing Forum 2", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Fullmetal"},
	{ID: 3, Title: "Testing Forum 3", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Naruto"},
}

// API definition

// Retrieve all the comments
func getComments(w http.ResponseWriter, r *http.Request) {
	var comments = Comment{}
	var comments_slice = []Comment{}

	db := dbConn()
	row, _ := db.Query("SELECT * FROM comments")

	for row.Next() {
		var title, text, date, author, anime string
		var id int64

		if err := row.Scan(&id, &title, &text, &author, &date, &anime); err != nil {
			panic(err.Error())
		}

		comments.ID = id
		comments.Title = title
		comments.Text = text
		comments.Author = author
		comments.Date = date
		comments.Anime = anime

		comments_slice = append(comments_slice, comments)
	}

	var response = JsonResponseComments{Type: "success", Data: comments_slice}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}

// Retrieve single comment by ID
func getCommentByID(w http.ResponseWriter, r *http.Request) {
	var comments = Comment{}
	var comments_slice = []Comment{}
	var title, text, date, author, anime string
	var row_id int64

	// Get the ID passed as a param in the url parsing the URL
	id, _ := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/comment/"), 10, 64)

	db := dbConn()
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&id, &title, &text, &author, &date, &anime); err != nil {
		if err == sql.ErrNoRows {
			panic(err.Error())
		}
		panic(err.Error())
	}

	comments.ID = row_id
	comments.Title = title
	comments.Text = text
	comments.Author = author
	comments.Date = date
	comments.Anime = anime

	comments_slice = append(comments_slice, comments)

	var response = JsonResponseComments{Type: "success", Data: comments_slice}
	json.NewEncoder(w).Encode(response)

	defer db.Close()

}

// Post new comment
func postComment(w http.ResponseWriter, r *http.Request) {
	var newComment Comment
	decoder := json.NewDecoder(r.Body).Decode(&newComment)
	// Using the BindJSON function to the newComment, which is a struct, we can attach the data from the Json to the struct data.
	if err := decoder; err != nil {
		return
	}

	// Add the newComment to the comments slice already created
	comments = append(comments, newComment)

	// Send back a status "created" with the newComment json file
	var response = JsonResponseComments{Type: "success", Data: comments}
	json.NewEncoder(w).Encode(response)
}

// Delete comment
func deleteComment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/comment/delete/")

	db := dbConn()
	delForm, err := db.Prepare("DELETE FROM comments WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(id)

	var response = JsonResponseComments{Type: "success", Data: []Comment{}}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
}

// Update comment
func updateComment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/comment/update/")
	var newComment Comment

	decoder := json.NewDecoder(r.Body).Decode(&newComment)
	if err := decoder; err != nil {
		return
	}

	commentFound, index := loopStructUpdate(id)
	if commentFound {
		if newComment.Title != "" {
			comments[index].Title = newComment.Title
		}
		if newComment.Text != "" {
			comments[index].Text = newComment.Text
		}

		var response = JsonResponseComments{Type: "success", Data: comments}
		json.NewEncoder(w).Encode(response)
		return
	}

	var response = JsonResponseComments{Type: "not found", Data: comments}
	json.NewEncoder(w).Encode(response)
}

// Utility functions
func removeElement(slice []Comment, index int) []Comment {
	return append(slice[:index], slice[index+1:]...)
}

func loopStructUpdate(id string) (bool, int) {
	for index, comment := range comments {
		if strconv.FormatInt(comment.ID, 10) == id {
			return true, index
		}
	}

	return false, 0
}

func dbConn() (db *sql.DB) {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "comments",
		AllowNativePasswords: true,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return
}
