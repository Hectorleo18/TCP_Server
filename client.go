package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"io"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		fmt.Println(msg)

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
				id:     CMD_MSG,
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

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("msg > " + msg + "\n"))
}

func (c *client) file(file *os.File, nickname string, filename string) {
	//var count int64
	for {
		buf := make([]byte, 2048)
		//Read file content
		n, err := file.Read(buf)
		if err != nil && io.EOF == err {
			fmt.Println("File Transfer")
			//Tell the server end file reception
			c.conn.Write([]byte("finish"))
			return
		}
		//Send to the server
		c.conn.Write([]byte(nickname+" ha enviado " + filename + "\n"))
		c.conn.Write(buf[:n])
	}
	// c.conn.Write()
}