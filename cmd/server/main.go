package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/slugalisk/shell/proto/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ServerOptions ...
type ServerOptions struct {
	Host     string
	Port     int
	CertPath string
	KeyPath  string
}

var options ServerOptions

func init() {
	flag.StringVar(&options.Host, "host", "localhost", "command server host")
	flag.IntVar(&options.Port, "port", 30013, "command server tcp port")
	flag.StringVar(&options.CertPath, "cert-path", "../../certs/server.pem", "tls certificate")
	flag.StringVar(&options.KeyPath, "key-path", "../../certs/server-key.pem", "tls key")
}

// Server ...
type Server struct {
	server *grpc.Server
}

// Serve ...
func (s *Server) Serve(options ServerOptions) error {
	creds, err := credentials.NewServerTLSFromFile(options.CertPath, options.KeyPath)
	if err != nil {
		return fmt.Errorf("could not load TLS keys: %s", err)
	}

	s.server = grpc.NewServer(grpc.Creds(creds))
	shell.RegisterShellServer(s.server, NewService())

	l, err := net.Listen("tcp", net.JoinHostPort(options.Host, strconv.Itoa(options.Port)))
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

func main() {
	server := Server{}
	go log.Fatal(server.Serve(options))

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	server.Stop()
}
