package main

import (
	"go-blueprints/trace"
	"log"
	"net/http"

	"github.com/stretchr/objx"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incoming messages that should be forwarded to the other clients.
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

// Note on pointers vs values: It makes no sense to pass by value because the channels in a room should be persistent
// entities, and the room itself is pretty much a global object.
// Concurrent access to the channels should be covered by channel mechanics anyway.
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("Client joined: ", client.userData["name"])
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left: ", client.userData["name"])
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// I guess each incoming client will spawn a new request handler, upgrading to a websocket until it's disconnected.
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// upgrade incoming HTTP request to WebSocket connection
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")

	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	// Create client, pass it a reference to the room so read() can forward incoming messages
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client

	// assumption: this is called when websocket connection is lost, which causes client.read to unblock and the
	// function to end.
	// Might be neater to have the client (on error) remove itself from the room? Would have two points where
	// we'd have to leave, while now we could add additional listeners for e.g. websocket events (if those are a thing).
	defer func() {
		r.leave <- client
	}()

	// both of these trigger infinite loops until the websocket borks.
	// write is called as a goroutine so it does not block and client.read is started
	go client.write()
	// client.read will keep this handler (instance?) open / spinning indefinitely until the websocket conks out.
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
