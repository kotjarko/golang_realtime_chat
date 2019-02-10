package main

import (
	"log"
	"time"
)

type Message struct {
	message_id int64
	from_user int64
	to_user int64
	text string
	dt time.Time
}

func newMessage(from_user int64, to_user int64, text string) *Message {
	result, err := DBConnection.Exec("INSERT INTO chat(from_user, to_user, text) VALUES($1, $2, $3)", from_user, to_user, text)

	if err != nil {
		log.Fatal("Cant create message: " + err.Error())
	}

	id, _ := result.LastInsertId()

	return &Message{id, from_user, to_user, text, time.Now()}
}


func listMessage(from_user int64, to_user int64) []*Message {
	// TODO add limit or show only from last X days

	rows, err := DBConnection.Query("SELECT * FROM chat " +
		"WHERE (from_user = $1 AND to_user = $2) or (from_user = $2 AND to_user = $1) ORDER BY dt DESC",
		from_user, to_user)

	if err != nil {
		log.Fatal("Cant get messages list")
	}
	defer rows.Close()

	list := make([]*Message, 0)
	for rows.Next() {
		ms := new(Message)
		err := rows.Scan(&ms.message_id, &ms.from_user, &ms.to_user, &ms.text, &ms.dt)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, ms)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return list
}
