package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/valyala/gorpc"
)

type Client struct {
	ID      string   `short:"s" long:"server" description:"id of the server to send the message to (from the configuration file) -- if empty, will be forwarded to all servers"`
	Pings   []string `short:"p" long:"ping" description:"optional user to ping -- supports '@' for ops+, and '*' for all users"`
	Targets []string `short:"c" long:"channel" description:"channel to send message to -- supports '*' for all joined channels" required:"true"`
}

func (c *Client) Usage() string {
	return "\"your example message here\""
}

func (c *Client) Execute(text []string) error {
	rpc := gorpc.NewUnixClient(conf.SocketFile)
	rpc.LogError = gorpc.NilErrorLogger
	rpc.Start()
	defer rpc.Stop()

	dp := newRpc()
	dc := dp.NewServiceClient("Daemon", rpc)

	if len(text) == 0 {
		return errors.New("no message specified (see 'client -h' for details')")
	}

	e := &Event{
		ID:      c.ID,
		Pings:   c.Pings,
		Targets: c.Targets,
		Text:    text,
	}

	if e.Pings == nil {
		e.Pings = []string{}
	}

	if len(e.Text) == 1 && strings.Contains(e.Text[0], "\n") {
		e.Text = strings.Split(e.Text[0], "\n")
	}

	resp, err := dc.CallTimeout("Send", e, 3*time.Second)
	if err != nil {
		cerr, ok := err.(*gorpc.ClientError)
		if ok && cerr.Timeout {
			return errors.New("rpc: timed out while sending request (is the daemon running?)")
		}

		return fmt.Errorf("rpc: %s", err)
	}

	fmt.Println(resp)

	return nil
}
