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

// var users map[string][]byte = make(map[string][]byte)
// var idxUsers int = 0

//getTokenUserPassword returns a jwt token for a user if the //password is ok
func getTokenUserPassword(w http.ResponseWriter, r *http.Request) {
	var login Login
	var user Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	//here I have a user!
	//Now check if exists
	db := dbConn()
	row := db.QueryRow("SELECT * FROM user WHERE username = ?", login.Username)
	defer db.Close()
	if err := row.Scan(&user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			println("Row no found.")
			log.Fatal(err)
		}
		println("Couldn't find searched params.")
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		return
	}
	token, err := createToken(login.Username)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var login Login
	var exist bool
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Cannot decode request", http.StatusBadRequest)
		return
	}
	db := dbConn()
	row := db.QueryRow("SELECT * FROM user WHERE username = ?", login.Username)
	defer db.Close()
	if err := row.Scan(&exist); err != nil && err != sql.ErrNoRows {
		println("Error in the query.")
		log.Fatal(err)
	} else if exist == true {
		println("User already registered.")
	}
	//If I'm here-> add user and return a token
	value, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)

	prepared, err_preparing := db.Prepare("INSERT INTO users(username,password) VALUES(?,?)")
	defer db.Close()
	if err_preparing != nil {
		println("Error preparing the mysql command.")
		log.Fatal(err_preparing)
	}

	res, err := prepared.Exec(login.Username, value)
	if err != nil {
		println("Error executing the mysql insertion.")
		log.Fatal(err)
		log.Fatal(res)
	}

	token, err := createToken(login.Username)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func createToken(username string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	secret := getSecret()
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
