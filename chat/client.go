package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan *message
	room   *room
	// TODO: oauth response, can probably be mapped to a struct.
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	// wait for incoming message; if an error is thrown, assumes disconnect? and exits loop.
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		// enrich incoming message with timestamp and poster
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	// range on channel, guess it waits forever?
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
