package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	defer c.socket.Close()
	// wait for incoming message; if an error is thrown, assumes disconnect? and exits loop.
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	// range on channel, guess it waits forever?
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
