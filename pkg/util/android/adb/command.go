package adb

import (
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

const adbExec = "adb"

type Cmd struct {
	ctx     context.Context
	cancel  func()
	device  string
	root    bool
	shell   bool
	buffer  io.Writer
	host    string
	command []string
	verbose bool
}

func NewCommand(args ...string) *Cmd {
	return &Cmd{
		command: args,
		ctx:     context.Background(),
		cancel:  func() {}, // empty cancel func that will get overridden if timeout is set
	}
}

func (c *Cmd) Context() context.Context {
	return c.ctx
}

func (c *Cmd) Cancel() {
	c.cancel()
}

func (c *Cmd) WithHost(host string) *Cmd {
	c.host = host
	return c
}

func (c *Cmd) WithContext(ctx context.Context) *Cmd {
	c.ctx = ctx
	return c
}

func (c *Cmd) WithTimeout(dur time.Duration) *Cmd {
	c.ctx, c.cancel = context.WithTimeout(c.ctx, dur)
	return c
}

func (c *Cmd) WithRoot() *Cmd {
	c.root = true
	return c
}

func (c *Cmd) WithDevice(dev string) *Cmd {
	c.device = dev
	return c
}

func (c *Cmd) WithShell() *Cmd {
	c.shell = true
	return c
}

func (c *Cmd) WithBuffer(writer io.Writer) *Cmd {
	c.buffer = writer
	return c
}

func (c *Cmd) WithVerbose() *Cmd {
	c.verbose = true
	return c
}

func (c *Cmd) Execute() ([]byte, error) {
	// defer cancel context to prevent leaks in case a timeout was set
	defer c.cancel()
	cmdArgs := make([]string, 0)
	if c.device != "" {
		cmdArgs = append(cmdArgs, "-s", c.device)
	}
	if c.host != "" {
		cmdArgs = append(cmdArgs, "-H", c.host)
	}
	if c.shell {
		cmdArgs = append(cmdArgs, "shell")
		if c.root {
			cmdArgs = append(cmdArgs, "su", "root")
		}
	}
	cmdArgs = append(cmdArgs, c.command...)

	var out []byte
	var err error
	if c.verbose {
		log.Println(append([]string{adbExec}, cmdArgs...))
	}
	if c.buffer != nil {
		cmd := exec.CommandContext(c.Context(), adbExec, cmdArgs...)
		cmd.Stdout = c.buffer
		err = cmd.Run()
	} else {
		out, err = exec.CommandContext(c.Context(), adbExec, cmdArgs...).Output()
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("%s: %s", exitErr.Error(), string(exitErr.Stderr))
		}
	}
	return out, err
}
