package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/slugalisk/shell/proto/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var certPath string
var host string
var port int
var timeout int64

func init() {
	flag.StringVar(&certPath, "cert-path", "../../certs/server.pem", "server certificate path")
	flag.StringVar(&host, "host", "localhost", "command server host")
	flag.IntVar(&port, "port", 30013, "command server tcp port")
	flag.Int64Var(&timeout, "timeout", 60, "command timeout")
}

// Client client wrapper
type Client struct {
	client shell.ShellClient
}

// NewClient create client wrapper
func NewClient(certPath string, host string, port int) (*Client, error) {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(
		net.JoinHostPort(host, strconv.Itoa(port)),
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client: shell.NewShellClient(conn),
	}
	return c, nil
}

// Exec ...
func (c *Client) Exec(timeout int64, name string, args ...string) error {
	client, err := c.client.Exec(
		context.Background(),
		&shell.ExecRequest{
			Command: &shell.Command{
				Name:    name,
				Args:    args,
				Timeout: timeout,
			},
		},
	)
	if err != nil {
		return err
	}

	for {
		rs, err := client.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("%s: %s", rs.Output.ShellId, rs.Output.Line)
	}
}

func main() {
	flag.Parse()

	// read command line arguments after --
	var command []string
	for i, v := range os.Args {
		if v == "--" {
			command = os.Args[i+1:]
			break
		}
	}

	if len(command) == 0 {
		log.Fatal("missing command")
	}

	client, err := NewClient(certPath, host, port)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Exec(timeout, command[0], command[1:]...); err != nil {
		log.Fatal(err)
	}
	log.Println("done")
}
