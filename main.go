package main

import (
	"log"
	"net"
)

func main() {
	//Crea un servidor
	s := newServer()
	go s.run()
	//Empieza a escuchar el puerto por el que recibirá información
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")
	//Acepta las conecciones de los clientes
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
