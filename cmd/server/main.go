package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/slugalisk/shell/proto/go"

	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	service *Service
	server  *grpc.Server
}

// NewServer ...
func NewServer() *Server {
	s := &Server{
		service: NewService(),
		server:  grpc.NewServer(),
	}
	shell.RegisterShellServer(s.server, s.service)

	return s
}

// Serve ...
func (s *Server) Serve(host string, port int) error {
	l, err := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	log.Printf("Serving at %s", l.Addr())
	return s.server.Serve(l)
}

// Stop ...
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}

var host string
var port int

func init() {
	flag.StringVar(&host, "host", "localhost", "command server host")
	flag.IntVar(&port, "port", 30013, "command server tcp port")
}

func main() {
	server := NewServer()

	go server.Serve(host, port)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	server.Stop()
}
