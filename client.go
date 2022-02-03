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
	for {
		//Read the channel until find a line break
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}
		//Remove the line break
		msg = strings.Trim(msg, "\r\n")
		//Separate the string in []string
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		//What command is?
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

//err send a message with a error
func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

//msg send a text message to the server
func (c *client) msg(msg string, nick string) {
	c.conn.Write([]byte("msg " + nick + " > " + msg + "\n"))
}

//file send a file to the server
func (c *client) file(file string, nickname string, filename string) {
	c.conn.Write([]byte("file " + nickname + " envió " + filename + " " + file + "\n"))
}
