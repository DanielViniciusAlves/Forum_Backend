package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// Main function

func main() {
	// Initialize the gin router
	router := gin.Default()

	// Associate the endpoint "/comments" with the getComments function.
	router.GET("/comments", getComments)

	// Associate the endpoint "/comments" with the getCommentByID function.
	router.GET("/comments/:id", getCommentByID)

	// Associate the endpoint "/comments" with the postComment function.
	router.POST("/comments", postComment)

	// Associate the endpoint "/comments" with the deleteComment function.
	router.DELETE("/comments/delete/:id", deleteComment)

	// Associate the endpoint "/comments" with the deleteComment function.
	router.PUT("/comments/update/:id", updateComment)

	// Start the associate the router whit one http server
	router.Run("localhost:8080")
}

// Seeding comments data

var comments = []comment{
	{ID: "1", Title: "Testing Forum", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Darling"},
	{ID: "2", Title: "Testing Forum 2", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Fullmetal"},
	{ID: "3", Title: "Testing Forum 3", Text: "Testing text", Author: "Daniel", Date: "28/11/2022", Anime: "Naruto"},
}

// API definition

// Retrieve all the comments
func getComments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, comments)
}

// Retrieve single comment by ID
func getCommentByID(c *gin.Context) {
	// Get the ID passed as a param in the url
	id := c.Param("id")

	// Loop the struct looking for the ID passed
	for _, comment := range comments {
		if comment.ID == id {
			c.IndentedJSON(http.StatusOK, comment)
			return
		}
	}

	// Will send NotFound status if the loop don't find the ID
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "comment not found"})
}

// Post new comment
func postComment(c *gin.Context) {
	var newComment comment

	// Using the BindJSON function to the newComment, which is a struct, we can attach the data from the Json to the struct data.
	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	// Add the newComment to the comments slice already created
	comments = append(comments, newComment)

	// Send back a status "created" with the newComment json file
	c.IndentedJSON(http.StatusCreated, newComment)
}

// Delete comment
func deleteComment(c *gin.Context) {
	id := c.Param("id")

	for index, comment := range comments {
		if comment.ID == id {
			comments = removeElement(comments, index)
			c.IndentedJSON(http.StatusOK, comments)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment could not be found."})
}

// Update comment
func updateComment(c *gin.Context) {
	id := c.Param("id")
	var newComment comment

	// fmt.Print("Test")
	if err := c.BindJSON(&newComment); err != nil {
		return
	} else if newComment.Title != "" || newComment.Text != "" {
		for index, comment := range comments {
			if comment.ID == id {
				if newComment.Title != "" {
					comments[index].Title = newComment.Title
				}
				if newComment.Text != "" {
					comments[index].Text = newComment.Text
				}
				c.IndentedJSON(http.StatusOK, comments[index])
				return
			}
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, newComment)
	}

}

// func updateCommentText(c *gin.Context){
// 	id :=
// }

// Utility functions
func removeElement(slice []comment, index int) []comment {
	return append(slice[:index], slice[index+1:]...)
}
