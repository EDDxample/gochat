package main

import (
	"flag"

	"github.com/EDDxample/gochat/pkg/client"
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "port")
)

func main() {
	flag.Parse()
	client := client.NewClient(*host, *port)
	client.Run()
}
