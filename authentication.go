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

//getTokenUserPassword returns a jwt token for a user if the password is ok
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
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", login.Username)
	defer db.Close()
	if err := row.Scan(&user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			println("Row not found.")
			log.Fatal(err)
		}
		println("Couldn't find searched params.")
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if err != nil {
		log.Fatal(err)
	}
	token, err := createToken(login.Username)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}
	// refresh, err := createRefreshToken(login.Username)
	// if err != nil {
	// 	http.Error(w, "Cannot create token", http.StatusInternalServerError)
	// 	return
	// }

	// var response Auth
	// response.Refresh = refresh
	// response.Token = token
	// print(response.Token)
	expirationTime := time.Now().Add(5 * time.Minute)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

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
	row := db.QueryRow("SELECT * FROM users WHERE username = ?", login.Username)
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
	// refresh, err := createRefreshToken(login.Username)
	// if err != nil {
	// 	http.Error(w, "Cannot create token", http.StatusInternalServerError)
	// 	return
	// }

	// var response Auth
	// response.Refresh = refresh
	// response.Token = token
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

func createRefreshToken(username string) (string, error) {
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

func refreshToken() {
	// 	Passive expiration

	// My favorite pattern is to have a dedicated server error.

	// Your server should respond with a particular error when the token is expired (to be distinguished from the 401 Unauthorized due to role access). You then add an HTTP middleware to your client that:

	//     detects this error response
	//     deletes local token and navigates to /auth/login

	// Or if you have a renew token:

	//     detects this error response
	//     attempts to renew the JWT
	//     repeats the original request on success OR navigates to auth page on failure.

	// This is a passive system that allows you to treat the JWT as an obscure string and does not have time-related issues.
	// 	"https://www.sohamkamani.com/golang/jwt-authentication/"
	// 	"https://docs.oracle.com/en/cloud/saas/live-experience/faled/handling-access-token-expiration.html"
}
