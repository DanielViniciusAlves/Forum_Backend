package main

// This file have the purpose of defining the endpoints.
import "net/http"

func endpointHandlers() {
	// Handling	get requests
	http.HandleFunc("/comments", getComments)
	http.HandleFunc("/comment/", getCommentByID)

	// Handling the post request for creating a new comment
	http.HandleFunc("/new_comment", postComment)

	// Handling delete request for specific comment
	http.HandleFunc("/comment/delete/", deleteComment)

	// Handling update request for specific comment
	http.HandleFunc("/comment/update/", updateComment)

	// Handling post request for creating new user
	http.HandleFunc("/new_user", createUser)

	// Handling post request for getting user token
	http.HandleFunc("/login", getTokenUserPassword)
}
