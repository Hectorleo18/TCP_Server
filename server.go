package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

/*
	Crea un nuevo servidor
*/
func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

/*
	Asigna la función a cada comando
*/
func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_FILE:
			s.msgFile(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		}
	}
}

/*
	Crea un nuevo cliente en el servidor
*/
func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

/*
	Cambia el nickname de un usuario
*/
func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME","server")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick), "server")
}

/*
	Agrega un cliente a una sala
*/
func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME", "server")
		return
	}

	roomName := args[1]
	//Revisa si la sala ya existe
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c
	//Elimina al usuario de la sala anterior si es que existe
	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName), "server")
}

/*
	Retorna todas las salas que han sido creadas
*/
func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")), "server")
}

/*
	Envía un mensaje de texto a los miembros de una sala
*/
func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG", "server")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, msg)
}

/*
	Envía un archivo a todos los miembros de una sala
*/
func (s *server) msgFile(c *client, args []string){
	if len(args) < 2 {
		c.msg("file is required, usage: /file FILE_PATH", "server")
		return
	}
	c.room.broadcastFile(c, args)
}

/*
	Termina la conexión con el cliente
*/
func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(", "server")
	c.conn.Close()
}

/*
	Elimina un usuario de una sala
*/
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
