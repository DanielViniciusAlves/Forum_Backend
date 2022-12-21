package main

// This file have the purpose of creating the middleware
import (
	"net/http"
)

func checkRequestToken(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validateJWTToken()
		originalHandler.ServeHTTP(w, r)
	})
}

func validateJWTToken() {

}
