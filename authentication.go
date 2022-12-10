package main

import "os"

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "testing"
	}
	return secret
}
