package server

import (
	"bufio"
	"fmt"
	"net"
)

// Server implementation for the TCP socket chat
type Server struct {
	host           string
	port           int
	chatMessages   chan string
	joiningClients chan chan string
	leavingClients chan chan string
}

// NewServer returns a new Server object
// and initializes the listener channels
func NewServer(host string, port int) *Server {
	return &Server{
		host:           host,
		port:           port,
		chatMessages:   make(chan string),
		joiningClients: make(chan chan string),
		leavingClients: make(chan chan string),
	}
}

// Run starts listening on the assigned port
// and bootstraps the message and connection handlers
func (s *Server) Run() {
	fmt.Printf("Listening on port %d...\n", s.port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		panic(err)
	}
	go s.HandleUserMessages()

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("error", err.Error())
			continue
		}
		go s.HandleConnection(conn)
	}
}

// HandleUserMessages keeps track of joining/leaving clients
// and broadcasts the new messages to all the active listeners
func (s *Server) HandleUserMessages() {
	activeClients := make(map[chan string]bool)

	for {
		select {

		// On client join
		case newClient := <-s.joiningClients:
			activeClients[newClient] = true

		// On client leave
		case leavingClient := <-s.leavingClients:
			delete(activeClients, leavingClient)
			close(leavingClient)

		// Broadcast Messages
		case message := <-s.chatMessages:
			for client := range activeClients {
				s.SendToClient(message, client)
			}

		}
	}
}

// HandleConnection adds a new client to the list
// of new joins and connects the input and output messages
// to their respective targets
func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	clientChannel := make(chan string)

	go s.ProcessOutboundMessages(conn, clientChannel)

	clientName := conn.RemoteAddr().String()

	// send join messages
	s.SendToClient(fmt.Sprintf("Welcome to the server, [%s]!\n", clientName), clientChannel)
	joinMsg := fmt.Sprintf("# [%s] joined the chat!\n", clientName)
	s.SendToServer(joinMsg)
	fmt.Printf(joinMsg)
	s.joiningClients <- clientChannel

	// scan for new client messages
	clientScanner := bufio.NewScanner(conn)
	for clientScanner.Scan() {
		s.SendToServer(fmt.Sprintf("[%s]: %s\n", clientName, clientScanner.Text()))
	}

	// send leave messages
	s.leavingClients <- clientChannel
	leaveMsg := fmt.Sprintf("# [%s] left the chat!\n", clientName)
	s.SendToServer(leaveMsg)
	fmt.Printf(leaveMsg)
}

// ProcessOutboundMessages writes the registered messages
// from the given channel to the client connection
func (s *Server) ProcessOutboundMessages(conn net.Conn, clientChannel <-chan string) {
	for message := range clientChannel {
		fmt.Fprintln(conn, message)
	}
}

// SendToServer registers the given text
// to be sent to the active clients
func (s *Server) SendToServer(text string) {
	s.chatMessages <- text
}

// SendToClient registers the given text
// to be sent to the given client channel
func (s *Server) SendToClient(text string, clientChannel chan string) {
	clientChannel <- text
}
