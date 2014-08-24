package gorc

import (
	"bufio"
	"net"
	"strings"
	"time"
)

var sesid int = 0

type Client struct {
	sesId    string
	name     string
	conn     net.Conn
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (client *Client) Read() {
	for {
		input, _ := client.reader.ReadString('\n')
		data := strings.Replace(input, "\r\n", "", 1)
		if len(data) <= 0 {
			if len(client.name) <= 0 {
				client.outgoing <- "Enter Your name: "
			}
			continue
		}

		if len(client.name) <= 0 {
			client.name = data
			continue
		}

		client.incoming <- NewProtocol(map[string]interface{}{
			"From":  client.name,
			"Body":  data,
			"SesId": client.sesId,
		}).JsonStringify()
	}
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func (client *Client) OnLogin() {
	client.outgoing <- "Hello! gorc\n"
	client.outgoing <- "Today is " + time.Now().Format(time.ANSIC) + "\n\n\n"
	client.outgoing <- "Enter Your name: "
}

func CreateClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	sesid++

	client := &Client{
		sesId:    string(sesid),
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
		conn:     connection,
	}

	client.Listen()
	client.OnLogin()

	return client
}
