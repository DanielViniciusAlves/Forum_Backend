# Defining API Endpoints
# This file need to be updated!
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

# Defining Project Architecture

Because this is a basic project we will use a flat structure to organize out project. With a main file, migrations and models for start. As the project evolves the idea is to add more files and when the structure becomes to much messy we can evolve for a more sophisticate architecture.

# Deciding Authentication

So I decided to using JWT to generate the token. The way Its gonna work is that once the user login in the application a Token and a temporary token will be generated, once the temporary token is expired the client side will refresh the Temporary Token with the one with a longer life spam. In this request a verification of the last time the user made a login will be made so the server can decide on asking for a new login or not.