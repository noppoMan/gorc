package gorc

import "net"

type Server struct {
	clients  []*Client
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (server *Server) Broadcast(data string) {
	for _, client := range server.clients {
		parsed := parseStreaam(data)
		if client.name == parsed.from {
			continue
		}
		//send message
		client.outgoing <- parsed.body
	}
}

func (server *Server) Join(connection net.Conn) {
	client := CreateClient(connection)
	server.clients = append(server.clients, client)
	go func() {
		for {
			server.incoming <- <-client.incoming
		}
	}()
}

func (server *Server) Listen() {
	go func() {
		for {
			select {
			case data := <-server.incoming:
				server.Broadcast(data)
			case conn := <-server.joins:
				server.Join(conn)
			}
		}
	}()
}

func CreateServer(port string) *Server {
	server := &Server{
		clients:  make([]*Client, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan string),
		outgoing: make(chan string),
	}

	server.Listen()

	listener, _ := net.Listen("tcp", ":"+port)
	println("gorc listening on tcp:" + port)

	for {
		conn, _ := listener.Accept()
		server.joins <- conn
	}
}
