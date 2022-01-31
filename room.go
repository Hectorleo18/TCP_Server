package main

import (
	"net"
	"strings"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

/*
	Envía un mensaje de texto a todos los integrantes de una sala
*/
func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg, sender.nick)
		}
	}
}

/*
	Envía un archivo a todos los integrantes de una sala
*/
func (r *room) broadcastFile(sender *client, args []string) {
	text := strings.Join(args[2:], " ")
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.file(text, sender.nick, args[1])
		}
	}
}
