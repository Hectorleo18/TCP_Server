package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	//Loop infinito
	for {
		//Lee lo que esté en el canal hasta que encuentre un salto de línea
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		//Se elimina el salto de línea
		msg = strings.Trim(msg, "\r\n")
		//Se separa la cadena por palabras
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		//Se comprueba qué comando fue el que se ingresó
		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/file":
			c.commands <- command{
				id:     CMD_FILE,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

/*
	Envía un mensaje notificando un error
*/
func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

/*
	Envía un mensaje de texto
*/
func (c *client) msg(msg string, nick string) {
	c.conn.Write([]byte("msg "+nick+" > " + msg + "\n"))
}

/*
	Envía un archivo
*/
func (c *client) file(file string, nickname string, filename string) {
	c.conn.Write([]byte("file "+nickname+" envió "+filename+" "+file+"\n"))
}