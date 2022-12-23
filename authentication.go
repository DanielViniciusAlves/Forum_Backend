package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "It should be panicking"
	}
	return secret
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var login, user Login

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}

	db := dbConn()
	if err := db.QueryRow("SELECT * FROM users WHERE username = ?", login.Username).Scan(&user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Row not found.", http.StatusBadRequest)
			println("Row not found.")
			log.Fatal(err)
		}
		println("Couldn't find searched params.")
		log.Fatal(err)
	}
	defer db.Close()

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		http.Error(w, "Wrong password.", http.StatusBadRequest)
		log.Fatal(err)
	}

	accessToken, err := createToken(login.Username, 15)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	// refreshToken, err := createToken(login.Username, 60)
	// if err != nil {
	// 	http.Error(w, "Cannot create refresh token", http.StatusInternalServerError)
	// 	return
	// }

	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "refresh_token",
	// 	Value: refreshToken,
	// })

	w.Header().Set("Set-Cookie", "first-cookie=value1")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var login Login
	var exist bool

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Cannot decode request", http.StatusBadRequest)
		return
	}
	db := dbConn()
	if err := db.QueryRow("SELECT * FROM users WHERE username = ?", login.Username).Scan(&exist); err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error in the query.", http.StatusInternalServerError)
		println("Error in the query.")
		log.Fatal(err)
	} else if exist {
		http.Error(w, "User already registered.", http.StatusBadRequest)
		println("User already registered.")
	}
	defer db.Close()

	value, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error generating hash password.", http.StatusInternalServerError)
		println("Error generating hash password.")
	}

	if prepared, err := db.Prepare("INSERT INTO users(username,password) VALUES(?,?)"); err != nil {
		http.Error(w, "Error in the query.", http.StatusInternalServerError)
		println("Error preparing the mysql command.")
	} else {
		if _, err := prepared.Exec(login.Username, value); err != nil {
			http.Error(w, "Error in the Execution.", http.StatusInternalServerError)
			println("Error executing the mysql insertion.")
		}
	}

	accessToken, err := createToken(login.Username, 15)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := createToken(login.Username, 60)
	if err != nil {
		http.Error(w, "Cannot create refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

func createToken(username string, expirationTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"username":   username,
		"exp":        time.Now().Add(time.Minute * expirationTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(getSecret()))
	if err != nil {
		return "Error signing token.", err
	}

	return tokenString, nil
}
