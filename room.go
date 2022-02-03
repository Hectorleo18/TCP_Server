package main

import (
	"net"
	"strings"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

//broadcast send a text message to all the clients in the room
func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg, sender.nick)
		}
	}
}

//broadcastFile send a file to all the members of the room
func (r *room) broadcastFile(sender *client, args []string) {
	text := strings.Join(args[2:], " ")
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.file(text, sender.nick, args[1])
		}
	}
}
