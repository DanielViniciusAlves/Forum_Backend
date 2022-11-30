package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	migrate()
	endpointHandlers()
	fmt.Println("Server listening in port 4000")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal(err)
	}
}
