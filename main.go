package main

import (
	"log"
	"net"
)

func main() {
	//Create a new server
	s := newServer()
	go s.run()
	//Start listen the port 8888
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")
	//Accept the conections of the clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
