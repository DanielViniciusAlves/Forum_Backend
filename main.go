package main

import (
	"log"
	"net/http"
)

func main() {
	migrate()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	endpointHandlers()
	println("Server listening in port 4000")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		println("Error while starting the server.")
		log.Fatal(err)
	}
}
