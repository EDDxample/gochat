package client

import (
	"fmt"
	"io"
	"net"
	"os"
)

// Client implementation for the TCP socket chat
type Client struct {
	host string
	port int
}

// NewClient returns a new Client object
func NewClient(host string, port int) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

// Run connects to the server and starts routing messages
func (c *Client) Run() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		panic(err)
	}

	done := make(chan bool)

	go func() {
		io.Copy(os.Stdout, conn)
		done <- true
	}()

	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		panic(err)
	}

	<-done

}
