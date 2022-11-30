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
	row := db.QueryRow("SELECT * FROM comments WHERE id = ?", id)
	if err := row.Scan(&row_id, &title, &text, &author, &date, &anime); err != nil {
		if err == sql.ErrNoRows {
			panic(err.Error())
		}
		panic(err.Error())
	}

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

	db := dbConn()
	decoder := json.NewDecoder(r.Body).Decode(&newComment)
	// // Using the BindJSON function to the newComment, which is a struct, we can attach the data from the Json to the struct data.
	if err := decoder; err != nil {
		panic(err.Error())
	}
	row, err := db.Prepare("INSERT INTO comments(title, comment_text, author, publish_date, anime) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	row.Exec(newComment.Title, newComment.Text, newComment.Author, newComment.Date, newComment.Anime)

	// fmt.Printf("test %s", newComment.Text)
	// // Send back a status "created" with the newComment json file
	var response = JsonResponseComments{Type: "success", Data: []Comment{}}
	json.NewEncoder(w).Encode(response)

	defer db.Close()
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

	db := dbConn()
	if newComment.Title != "" {
		insForm, err := db.Prepare("UPDATE comments SET title=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(newComment.Title, id)
	}
	if newComment.Text != "" {
		insForm, err := db.Prepare("UPDATE comments SET comment_text=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(newComment.Text, id)
	}

	var response = JsonResponseComments{Type: "success", Data: []Comment{}}
	json.NewEncoder(w).Encode(response)
}

// Utility functions
func removeElement(slice []Comment, index int) []Comment {
	return append(slice[:index], slice[index+1:]...)
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
