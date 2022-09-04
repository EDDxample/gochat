package main

import (
	"flag"

	"github.com/EDDxample/gochat/pkg/server"
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "port")
)

func main() {
	flag.Parse()
	server := server.NewServer(*host, *port)
	server.Run()
}
