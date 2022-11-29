# Defining API Endpoints
* This document will have the purpose of defining the projects endpoints for the API.

First we will define the idea of the project. Basically it will be for starter a forum. The user will have the possibility of make a comment, edit this comment, delete and obviously see the comment. So it will be a CRUD.

The endpoint will be:

```

/comments

    Get -> It will retrieve a list of all the comments.
    Post -> It will be responsible for making a new comment.

/comments/:id

    Get -> Will retrieve a specific comment requested by the user.
    Update -> Will update a specific comment.

/comments/delete/:id

    Delete -> Will delete a specific comment.

```