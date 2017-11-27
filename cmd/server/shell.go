package main

import (
	"io"
	"sync"

	"github.com/satori/go.uuid"
	"github.com/slugalisk/shell/proto/go"
)

// Shell ...
type Shell struct {
	id           string
	server       shell.Shell_FollowServer
	commandsLock sync.Mutex
	commands     map[string]*Command
}

// NewShell ...
func NewShell(server shell.Shell_FollowServer) *Shell {
	return &Shell{
		id:       uuid.NewV4().String(),
		server:   server,
		commands: make(map[string]*Command),
	}
}

// ID ...
func (s *Shell) ID() string {
	return s.id
}

func (s *Shell) removeAllCommands() {
	s.commandsLock.Lock()
	defer s.commandsLock.Unlock()

	for id, command := range s.commands {
		command.HandleExit(&shell.CommandExit{
			CommandId: id,
			ShellId:   s.id,
			Code:      1,
		})
	}
}

func (s *Shell) removeCommand(exit *shell.CommandExit) {
	s.commandsLock.Lock()
	command, ok := s.commands[exit.CommandId]
	s.commandsLock.Unlock()

	if ok {
		exit.ShellId = s.id
		command.HandleExit(exit)

		s.commandsLock.Lock()
		delete(s.commands, exit.CommandId)
		s.commandsLock.Unlock()
	}
}

func (s *Shell) handleOutput(output *shell.CommandOutput) {
	s.commandsLock.Lock()
	command, ok := s.commands[output.CommandId]
	s.commandsLock.Unlock()

	if ok {
		output.ShellId = s.id
		command.HandleOutput(output)
	}
}

// HandleResponses ...
func (s *Shell) HandleResponses() error {
	defer s.removeAllCommands()

	for {
		res, err := s.server.Recv()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}

		if exit := res.GetExit(); exit != nil {
			s.removeCommand(exit)
		}
		if output := res.GetOutput(); output != nil {
			s.handleOutput(output)
		}
	}
}

// Exec ...
func (s *Shell) Exec(command *Command) {
	s.commandsLock.Lock()
	s.commands[command.GetId()] = command
	s.commandsLock.Unlock()

	// send message to command channel
	s.server.Send(&shell.FollowResponse{
		Command: command.Command,
	})
}
