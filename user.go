package main

import (
	"database/sql"
	"log"
)

/*
user_id serial primary key,
name varchar(30)
 */
type User struct {
	user_id int64
	name string
}

func newUser(name string) *User {
	// check existing user
	row := DBConnection.QueryRow("SELECT * FROM users WHERE name = $1", name)
	us := new(User)
	err := row.Scan(&us.user_id, &us.name)

	if err != sql.ErrNoRows {
		return us
	}

	// create user
	result, err := DBConnection.Exec("INSERT INTO users(name) VALUES($1)", name)

	if err != nil {
		log.Fatal("Cant create user: " + err.Error())
	}

	id, _ := result.LastInsertId()

	return &User{id, name}
}


func listUser(currentId int64) []*User {
	rows, err := DBConnection.Query("SELECT * FROM users WHERE user_id <> $1", currentId)
	if err != nil {
		log.Fatal("Cant get users list")
	}
	defer rows.Close()

	list := make([]*User, 0)
	for rows.Next() {
		us := new(User)
		err := rows.Scan(&us.user_id, &us.name)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, us)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return list
}

func getUser(id int64) *User {
	row := DBConnection.QueryRow("SELECT * FROM users WHERE user_id = $1", id)
	us := new(User)
	err := row.Scan(&us.user_id, &us.name)

	if err == sql.ErrNoRows {
		log.Fatal("User not found")
	}

	return us
}