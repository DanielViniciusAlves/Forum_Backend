package main

import (
	"database/sql"
	"log"
)

func migrate() {
	db := dbConn()
	createDataBase(db)
	populateDataBase(db)

	defer db.Close()
}

func createDataBase(db *sql.DB) {
	db.Exec("use comments")
	db.Exec("DROP TABLE IF EXISTS comments")
	_, err := db.Exec("CREATE TABLE comments (id INT AUTO_INCREMENT NOT NULL, title VARCHAR(128) NOT NULL, comment_text VARCHAR(255) NOT NULL, author VARCHAR(255) NOT NULL, publish_date VARCHAR(255) NOT NULL, anime VARCHAR(255) NOT NULL, PRIMARY KEY (`id`) )")
	if err != nil {
		println("Error while creating comments table.")
		log.Fatal(err)
	}
}

func populateDataBase(db *sql.DB) {
	_, err := db.Exec("INSERT INTO comments (title, comment_text, author, publish_date, anime) VALUES ('Testing 1', 'Im testing the mysql 1', 'Daniel', '29/11/2022', 'Darling'), ('Testing 2', 'Im testing the mysql 2', 'Daniel', '29/11/2022', 'Boruto'), ('Testing 3', 'Im testing the mysql 3', 'Daniel', '29/11/2022', 'Naruto'), ('Testing 4', 'Im testing the mysql 4', 'Daniel', '29/11/2022', 'Fullmetal')")
	if err != nil {
		println("Error while populating comments table.")
		log.Fatal(err)
	}

}
