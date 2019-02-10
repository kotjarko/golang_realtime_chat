package main

import "golang.org/x/net/websocket"


type Chat struct {
	from_user int64
	to_user int64
	messages []*Message
	ws *websocket.Conn
}

func newChat(from_user int64, to_user int64, conn *websocket.Conn) *Chat {
	return &Chat{from_user, to_user, listMessage(from_user, to_user), conn}
}

func (c *Chat) sendMessage(text string) *Message {
	message := newMessage(c.from_user, c.to_user, text)
	c.messages = append(c.messages, message)
	return message
}