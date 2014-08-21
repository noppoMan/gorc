package gorc

import (
	"bufio"
	"net"
	"strings"
	"time"

	"github.com/noppoman/gorc/util"
)

const NAME_PROMPT string = "Enter Your name: "
const SCHEMA string = "gorc://"

type Protocol struct {
	schema string
	from   string
	body   string
}

func parseStream(stream string) *Protocol {
	replaced := strings.Replace(stream, SCHEMA, "", 1)
	splited := strings.Split(replaced, ",")
	f := splited[0]
	b := splited[1]

	proto := new(Protocol)
	proto.schema = SCHEMA
	proto.from = strings.Split(f, ":")[1]
	proto.body = util.Base64Decode(strings.Split(b, ":")[1])

	return proto
}

func toSendableData(from string, input string) string {
	return SCHEMA +
		"from:" +
		from +
		",body:" +
		util.Base64Encode(strings.Join([]string{"Sent From:" + from, "Date:" + time.Now().Format(time.ANSIC), input}, "\n")) + "\n\n"
}

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
				client.outgoing <- NAME_PROMPT
			}
			continue
		}

		if len(client.name) <= 0 {
			client.name = strings.Replace(input, "\n", "", 1)
			continue
		}

		client.incoming <- toSendableData(client.name, input)
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
	client.outgoing <- NAME_PROMPT
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
