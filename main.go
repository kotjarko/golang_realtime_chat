package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
)

var staticFolder = "static"
var DBConnection *sql.DB
var chats []*Chat

func init() {
	var err error
	DBConnection, err = sql.Open("postgres", "postgres://test:test@localhost/test")
	if err != nil {
		log.Fatal(err)
	}

	if err = DBConnection.Ping(); err != nil {
		log.Fatal(err)
	}
}

func chatHandler(ws *websocket.Conn) {
	var err error

	for {
		var reply []byte

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Cant receive")
			break
		}

		var dat map[string] interface{}
		if err := json.Unmarshal(reply, &dat); err != nil {
			log.Fatal("not json data")
		}
		fmt.Println("Recieve: ", dat)

		var response string

		switch dat["action"] {
			case "auth":
				usr := newUser(dat["name"].(string))
				chats = append(chats, &Chat{usr.user_id,0,nil,ws})
				response = "AUTH OK"

			case "list":
				for i := range chats {
					if chats[i].ws == ws {
						list := listUser(chats[i].from_user)

						res := make(map[string]string)
						for j:= range list {
							res[strconv.FormatInt(list[j].user_id, 10)] = list[j].name
						}

						resJSON, err := json.Marshal(res)
						if err != nil {
							panic(err)
						}
						response = string(resJSON)
						break
					}
				}
			case "join":
				for i := range chats {
					if chats[i].ws == ws {
						chats[i].to_user, err = strconv.ParseInt(dat["user"].(string), 10, 64)
						if err != nil {
							panic(err)
						}

						type ChatListMessage struct {
							Author int64 `json:"author"`
							Text string `json:"text"`
						}

						type ChatList struct {
							From string `json:"from"`
							To string `json:"to"`
							To_id int64 `json:"to_id"`
							Messages map[string]ChatListMessage `json:"messages"`
						}

						res := ChatList{
							getUser(chats[i].from_user).name,
							getUser(chats[i].to_user).name,
							chats[i].from_user,
							make(map[string]ChatListMessage)}

						chats[i].messages = listMessage(chats[i].from_user, chats[i].to_user)

						for j:= range chats[i].messages {
							res.Messages[chats[i].messages[j].dt.String()] = ChatListMessage{
								chats[i].messages[j].from_user,
								chats[i].messages[j].text}
						}

						resJSON, err := json.Marshal(res)
						if err != nil {
							panic(err)
						}

						response = string(resJSON)
						fmt.Println(res)
						fmt.Println(response)
						break
					}
				}
			case "send":
				for i := range chats {
					if chats[i].ws == ws {
						if chats[i].from_user != 0 && chats[i].to_user != 0 {
							var message = chats[i].sendMessage(dat["text"].(string))
							response = "SENT"
							res := struct {
								Dt string `json:"dt"`
								Author string `json:"author"`
								Text string `json:"text"`
							}{message.dt.String(), getUser(message.from_user).name, message.text}

							resJSON, err := json.Marshal(res)
							if err != nil {
								panic(err)
							}

							response = string(resJSON)

							for j := range chats {
								if chats[j].from_user == chats[i].to_user {
									websocket.Message.Send(chats[j].ws, response)
									break
								}
							}
						}
						break
					}
				}
		}
		if len(response) == 0 {
			response = "ERROR"
		}
		websocket.Message.Send(ws, response)
	}
}

func main() {
	staticServer := http.FileServer(http.Dir(staticFolder))
	http.Handle("/", staticServer)

	http.Handle("/chat", websocket.Handler(chatHandler))

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln("fatal error: ", err.Error())
	}
}
