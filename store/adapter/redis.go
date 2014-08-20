package adapter

import (
	"bufio"
	"fmt"
	"net"
)

type Logger struct{}

func (p *Logger) info(cmd, resp string, e error) {
	fmt.Printf(
		"#####\t%v\nRESPO\t%v\nERROR\t%v\n",
		cmd,
		resp,
		e,
	)
}

type RedisTalker struct {
	host       string
	port       string
	logger     Logger
	connection net.Conn
}

func (self *RedisTalker) initialize(host string, port string) RedisTalker {
	self.logger = new(Logger)
	self.connection, _ = net.Dial("tcp", "localhost:6379")
	return self
}

func (self *RedisTalker) set(key string, value string) {
	resp = make([]byte, 1024)
	reader = bufio.NewReaderSize(self.connection, 1024)
	fmt.Fprintf(self.connection, "*3\r\n$3\r\nSET\r\n$7\r\n"+key+"\r\n$5\r\n"+value+"\r\n")
}

func main() {

	var conn net.Conn
	var reader *bufio.Reader
	var resp []byte
	var err error
	var length int
	var rerr error

	talker := new(RedisTalker).initialize

	resp = make([]byte, 1024)

	conn, _ = net.Dial("tcp", "localhost:6379")
	reader = bufio.NewReaderSize(conn, 1024)

	fmt.Fprintf(conn, "*3\r\n$3\r\nSET\r\n$7\r\ntainaka\r\n$5\r\nritsu\r\n")

	// "SET mysample002 true\r\n" の結果の表示
	length, rerr = reader.Read(resp)
	//logger.info("SET", string(resp), err)
	fmt.Printf("READ STRING %v %v\n", length, rerr)
}
