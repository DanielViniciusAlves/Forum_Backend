package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Comments data structure.

type comment struct {
	ID     string `json: "id"`
	Title  string `json: "title"`
	Text   string `json: "text"`
	Author string `json: "author"`
	Date   string `json: "date"`
	Anime  string `json: "anime"`
}

// Json file that will be send

type JsonResponseComments struct {
	Type    string    `json:"type"`
	Data    []comment `json:"data"`
	Message string    `json:"message"`
}

type JsonResponseComment struct {
	Type    string  `json:"type"`
	Data    comment `json:"data"`
	Message string  `json:"message"`
}

// Main function

func main() {
	// Defining directory for the router to get the templates
	// router.LoadHTMLGlob("templates/*")

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

var comments = []comment{
	{ID: "1", Title: "Testing Forum", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Darling"},
	{ID: "2", Title: "Testing Forum 2", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Fullmetal"},
	{ID: "3", Title: "Testing Forum 3", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Naruto"},
}

// API definition

// Retrieve all the comments
func getComments(w http.ResponseWriter, r *http.Request) {
	var response = JsonResponseComments{Type: "success", Data: comments}
	json.NewEncoder(w).Encode(response)
}

// Retrieve single comment by ID
func getCommentByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID passed as a param in the url parsing the URL
	id := strings.TrimPrefix(r.URL.Path, "/comment/")

	// Loop the struct looking for the ID passed
	for _, comment := range comments {
		if comment.ID == id {
			var response = JsonResponseComment{Type: "success", Data: comment}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Will send NotFound status if the loop don't find the ID
	var response = JsonResponseComments{Type: "not found", Data: comments}
	json.NewEncoder(w).Encode(response)
}

// Post new comment
func postComment(w http.ResponseWriter, r *http.Request) {
	var newComment comment
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

	for index, comment := range comments {
		if comment.ID == id {
			comments = removeElement(comments, index)
			var response = JsonResponseComments{Type: "success", Data: comments}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	var response = JsonResponseComments{Type: "not found", Data: comments}
	json.NewEncoder(w).Encode(response)
}

// Update comment
func updateComment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/comment/update/")
	var newComment comment

	decoder := json.NewDecoder(r.Body).Decode(&newComment)
	if err := decoder; err != nil {
		return
	}

	commentFound, index := loopStruct(id)
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
func removeElement(slice []comment, index int) []comment {
	return append(slice[:index], slice[index+1:]...)
}

func loopStruct(id string) (bool, int) {
	for index, comment := range comments {
		if comment.ID == id {
			return true, index
		}
	}

	return false, 0
}
