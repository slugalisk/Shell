package main

import (
	"github.com/satori/go.uuid"
	"github.com/slugalisk/shell/proto/go"
)

// CommandOutputHandler ...
type CommandOutputHandler func(output *shell.CommandOutput)

// CommandExitHandler ...
type CommandExitHandler func(output *shell.CommandExit)

// Command ...
type Command struct {
	*shell.Command
	OnOutput CommandOutputHandler
	OnExit   CommandExitHandler
}

// NewCommand ...
func NewCommand(command *shell.Command) *Command {
	command.Id = uuid.NewV4().String()
	return &Command{
		Command: command,
	}
}

// HandleOutput ...
func (c *Command) HandleOutput(output *shell.CommandOutput) {
	if c.OnOutput != nil {
		c.OnOutput(output)
	}
}

// HandleExit ...
func (c *Command) HandleExit(exit *shell.CommandExit) {
	if c.OnExit != nil {
		c.OnExit(exit)
	}
}
