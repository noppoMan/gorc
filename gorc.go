package main

import (
	"flag"

	"github.com/noppoman/gorc/gorc"
)

func main() {

	var (
		port string
	)
	flag.StringVar(&port, "port", "6667", "Listening port")

	gorc.CreateServer(port)
}
