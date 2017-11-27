package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/slugalisk/shell/proto/go"
)

// Service ...
type Service struct {
	shellsLock sync.Mutex
	shells     map[string]*Shell
}

// NewService ...
func NewService() *Service {
	return &Service{
		shells: make(map[string]*Shell),
	}
}

// Time ...
func (s *Service) Time(ctx context.Context, req *shell.TimeRequest) (*shell.TimeResponse, error) {
	time, _ := ptypes.TimestampProto(time.Now())
	return &shell.TimeResponse{Time: time}, nil
}

// Ping ...
func (s *Service) Ping(ctx context.Context, req *shell.PingRequest) (*shell.PingResponse, error) {
	return &shell.PingResponse{Data: req.Data}, nil
}

// Exec ...
func (s *Service) Exec(req *shell.ExecRequest, server shell.Shell_ExecServer) error {
	wg := sync.WaitGroup{}

	command := NewCommand(req.Command)
	command.OnOutput = func(output *shell.CommandOutput) {
		server.Send(&shell.ExecResponse{
			Output: output,
		})
	}
	command.OnExit = func(exit *shell.CommandExit) {
		wg.Done()
	}

	s.shellsLock.Lock()
	wg.Add(len(s.shells))
	for _, shell := range s.shells {
		shell.Exec(command)
	}
	s.shellsLock.Unlock()

	wg.Wait()

	return nil
}

// Follow ...
func (s *Service) Follow(follower shell.Shell_FollowServer) error {
	shell := NewShell(follower)

	s.shellsLock.Lock()
	s.shells[shell.ID()] = shell
	s.shellsLock.Unlock()

	log.Printf("shell %s connected", shell.ID())
	log.Println(shell.HandleResponses())
	log.Printf("shell %s disconnected", shell.ID())

	s.shellsLock.Lock()
	delete(s.shells, shell.ID())
	s.shellsLock.Unlock()

	return nil
}
