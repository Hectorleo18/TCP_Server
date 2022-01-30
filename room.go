package main

import (
	"net"
	"os"
	"strings"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}

func (r *room) broadcastFile(sender *client, args []string) {
	// nickname := sender.nick
	// buf := make([]byte, 2048)
	// n, err := sender.conn.Read(buf)
	// if err != nil {
	// 	return
	// }
	// fileName := string(buf[:n])
	f, err := os.Create(args[1])
	if err != nil {
		return
	}
	defer f.Close()
	text := strings.Join(args[2:], " ")
	text = strings.ReplaceAll(text, "HectorLeoRodriguez", "\n")
	srcFile := []byte(text)
	f.Write(srcFile)
	// for {
	// 	buf := make([]byte, 2048)
	// 	n, _ := sender.conn.Read(buf)
	// 	if string(buf[:n]) == "finish" {
	// 		break
	// 	}
	// 	f.Write(buf[:n])
	// }
	// fi, err := os.Open(fileName)
	// if err != nil {
	// 	return
	// }
	// defer fi.Close()

	// for addr, m := range r.members {
	// 	if sender.conn.RemoteAddr() != addr {
	// 		m.file(fi, nickname, fileName)
	// 	}
	// }
}
