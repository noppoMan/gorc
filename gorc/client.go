package gorc

import (
	"bufio"
	"net"
	"strings"
	"time"
)

type Client struct {
	name     string
	incoming chan string
	outgoing chan string
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (client *Client) Read() {
	for {
		input, _ := client.reader.ReadString('\n')
		if len(input) <= 2 {
			if len(client.name) <= 0 {
				client.outgoing <- "Enter Your name: "
			}
			continue
		}

		if len(client.name) <= 0 {
			client.name = strings.Replace(input, "\n", "", 1)
			continue
		}

		// TODO implement quit command.

		client.incoming <- new(Protocol).Initialize(map[string]interface{}{
			"From": client.name,
			"Body": input,
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

	client := &Client{
		incoming: make(chan string),
		outgoing: make(chan string),
		reader:   reader,
		writer:   writer,
	}

	client.Listen()
	client.OnLogin()

	return client
}
