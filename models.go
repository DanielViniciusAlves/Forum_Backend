package main

type Comment struct {
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Anime  string `json:"anime"`
}

type CommentTest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Anime  string `json:"anime"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	Refresh_Token string `json:"refresh_token"`
	Access_Token  string `json:"access_token"`
}
