package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oyvindsk/go-blueprints-chat/trace"
	"github.com/stretchr/objx"
)

type room struct {
	// forward is the channel that holds messages that should be forwared to the other clients
	forward chan *message

	// a channel for ppl whishing to join the room
	join chan *client

	// a channel for ppl wanting to leave the room
	leave chan *client

	// all clients currently in the room
	clients map[*client]bool

	// will receive trace information about the activity on the room
	tracer trace.Tracer

	// avatar is how avatar information will be obtained
	avatar Avatar
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			// foreward message to all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send the message
					r.tracer.Trace(" -- sent to client")
				default:
					// failed to send
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- failed to send, cleaned up client")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	authCoookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCoookie.Value),
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}
