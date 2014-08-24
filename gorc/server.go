package gorc

import "net"

type Server struct {
	clients  []*Client
	joins    chan net.Conn
	incoming chan string
	outgoing chan string
}

func (server *Server) Broadcast(proto *Protocol) {
	for _, client := range server.clients {
		if client.name == proto.From {
			continue
		}
		//send message
		client.outgoing <- proto.Dispfy()
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
				proto := NewProtocolFromString(data)
				if proto.IsQuit() {
					for _, client := range server.clients {
						if client.sesId == proto.SesId {
							client.conn.Close()
						}
					}
					continue
				}
				server.Broadcast(proto)
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
