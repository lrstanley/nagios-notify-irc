// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/intel/tfortools"
	"github.com/lrstanley/girc"
	"github.com/valyala/gorpc"
)

func newRpc() *gorpc.Dispatcher {
	dp := gorpc.NewDispatcher()
	dp.AddService("Daemon", &Daemon{})

	return dp
}

type Daemon struct {
	StartupDelay int `short:"d" long:"delay" description:"delay (seconds) between each new irc connection" default:"1"`
}

func (s *Daemon) Ping() {}

func (s *Daemon) Send(event *Event) string {
	for i := 0; i < len(conf.Servers); i++ {
		if conf.Servers[i].ID == event.ID || event.ID == "" {
			conf.Servers[i].recv <- event
		}
	}

	return "OK"
}

func (s *Daemon) Execute([]string) error {
	done := make(chan struct{})
	var wg sync.WaitGroup

	if _, err := os.Stat(conf.SocketFile); err == nil {
		rpc := gorpc.NewUnixClient(conf.SocketFile)
		rpc.Start()

		dp := newRpc()
		dc := dp.NewServiceClient("Daemon", rpc)

		if _, err := dc.CallTimeout("Ping", nil, 1*time.Second); err != nil {
			fmt.Printf("removing stale socket file %q\n", conf.SocketFile)
			if err = os.Remove(conf.SocketFile); err != nil {
				rpc.Stop()
				return err
			}
		} else {
			rpc.Stop()
			return fmt.Errorf("error: daemon already found listening at %q", conf.SocketFile)
		}

		rpc.Stop()
	}

	dp := newRpc()
	rpc := gorpc.NewUnixServer(conf.SocketFile, dp.NewHandlerFunc())
	err := rpc.Start()
	if err != nil {
		return fmt.Errorf("rpc: %s", err)
	}

	for i := 0; i < len(conf.Servers); i++ {
		wg.Add(1)
		conf.Servers[i].recv = make(chan *Event)
		time.Sleep(time.Duration(s.StartupDelay) * time.Second)
		go conf.Servers[i].setup(done, &wg)
	}

	catch()
	rpc.Stop()
	close(done)
	wg.Wait()

	for i := 0; i < len(conf.Servers); i++ {
		close(conf.Servers[i].recv)
	}

	fmt.Println("exiting")
	return nil
}

type Event struct {
	ID      string   // The ID of the server from the configuration file.
	Pings   []string // "*", "@", or list of users.
	Targets []string // "*" or list of channels.
	Text    []string
}

func (e *Event) String() string {
	return fmt.Sprintf("<[Event] id:%q pings:%#v targets:%#v text:%#v>", e.ID, e.Pings, e.Targets, e.Text)
}

type Server struct {
	ID            string   `toml:"id"`
	Hostname      string   `toml:"hostname"`
	Password      string   `toml:"password"`
	Bind          string   `toml:"bind"`
	Port          int      `toml:"port"`
	TLS           bool     `toml:"tls"`
	TLSSkipVerify bool     `toml:"tls_skip_verify"`
	Channels      []string `toml:"channels"`
	DisableColors bool     `toml:"disable_colors"`
	Nick          string   `toml:"nick"`
	Name          string   `toml:"name"`
	User          string   `toml:"user"`
	SASLUser      string   `toml:"sasl_user"`
	SASLPass      string   `toml:"sasl_pass"`

	log  *log.Logger
	recv chan *Event
}

func (s *Server) String() string {
	return fmt.Sprintf(
		"<[Server] id:%q host:%q port:%d bind:%q tls:%t tls-skip-verify:%t nick:%q channels:%#v>",
		s.ID, s.Hostname, s.Port, s.Bind, s.TLS, s.TLSSkipVerify, s.Nick, s.Channels,
	)
}

func (s *Server) setup(done chan struct{}, wg *sync.WaitGroup) error {
	defer wg.Done()
	if s.ID == "" {
		return errors.New("empty server id specified")
	}

	s.log = log.New(os.Stdout, s.ID+": ", log.Ltime)

	if s.Port == 0 {
		s.Port = conf.DefaultPort
	}

	if s.Nick == "" {
		s.Nick = conf.DefaultNick
	}

	if s.Name == "" {
		s.Name = conf.DefaultName
	}

	if s.User == "" {
		s.User = conf.DefaultUser
	}

	s.log.Printf("adding %s", s.String())

	ircConfig := girc.Config{
		Server:       s.Hostname,
		ServerPass:   s.Password,
		Port:         s.Port,
		Nick:         s.Nick,
		Name:         s.Name,
		User:         s.User,
		Bind:         s.Bind,
		SSL:          s.TLS,
		Version:      "https://github.com/lrstanley/nagios-notify-irc " + version,
		GlobalFormat: !s.DisableColors,
		TLSConfig:    &tls.Config{ServerName: s.Hostname, InsecureSkipVerify: s.TLSSkipVerify},
		RecoverFunc:  func(_ *girc.Client, e *girc.HandlerError) { s.log.Print(e.Error()) },
	}

	if s.SASLUser != "" || s.SASLPass != "" {
		ircConfig.SASL = &girc.SASLPlain{User: s.SASLUser, Pass: s.SASLPass}
	}

	client := girc.New(ircConfig)
	client.Handlers.AddBg(girc.ALL_EVENTS, s.onAll)
	client.Handlers.Add(girc.CONNECTED, s.onConnect)

	var wgDone sync.WaitGroup
	go func() {
		wgDone.Add(1)
		for {
			err := client.Connect()
			if err == nil {
				break
			}

			s.log.Printf("error: %s", err)

			s.log.Printf("sleeping for %ds before reconnecting", conf.ReconnectDelay)
			time.Sleep(time.Duration(conf.ReconnectDelay) * time.Second)
		}

		wgDone.Done()
	}()

	for {
		select {
		case <-done:
			goto done
		case e := <-s.recv:
			s.handle(client, e)
		}
	}

done:
	client.Close()
	wgDone.Wait()

	return nil
}

func (s *Server) handle(c *girc.Client, e *Event) {
	s.log.Printf("handling event: %s", e)
	targets := []string{}
	for i := 0; i < len(e.Targets); i++ {
		if e.Targets[i] == "*" {
			targets = c.ChannelList()
			break
		}

		targets = append(targets, e.Targets[i])
	}

	var pingAll bool
	var pingOps bool
	for i := 0; i < len(e.Pings); i++ {
		if e.Pings[i] == "*" {
			pingAll = true
			break
		}

		if e.Pings[i] == "@" {
			pingOps = true
			break
		}
	}

	for i := 0; i < len(targets); i++ {
		channel := c.LookupChannel(targets[i])
		if channel == nil {
			s.log.Printf("requested send to unknown channel %q", targets[i])
			continue
		}

		if pingAll {
			c.Cmd.Message(targets[i], strings.Join(channel.UserList, " ")+":")
		} else if pingOps {
			users := channel.Admins(c)
			ops := []string{}
			for j := 0; j < len(users); j++ {
				ops = append(ops, users[j].Nick)
			}

			if len(ops) > 0 {
				c.Cmd.Message(targets[i], strings.Join(ops, " ")+":")
			}
		} else if len(e.Pings) > 0 {
			c.Cmd.Message(targets[i], strings.Join(e.Pings, " ")+":")
		}

		for j := 0; j < len(e.Text); j++ {
			if flags.Client.Plain {
				c.Cmd.Message(targets[i], e.Text[j])
				continue
			}

			buf := &bytes.Buffer{}
			if err := tfortools.OutputToTemplate(buf, "", e.Text[j], nil, nil); err != nil {
				s.log.Printf("error executing text template: %s", err)
				continue
			}

			c.Cmd.Message(targets[i], strings.TrimSpace(buf.String()))
		}
	}
}

func (s *Server) onConnect(c *girc.Client, e girc.Event) {
	for i := 0; i < len(s.Channels); i++ {
		if split := strings.SplitN(s.Channels[i], " ", 2); len(split) == 2 {
			c.Cmd.JoinKey(split[0], split[1])
			continue
		}

		c.Cmd.Join(s.Channels[i])
	}
}

func (s *Server) onAll(c *girc.Client, e girc.Event) {
	if out, ok := e.Pretty(); ok {
		s.log.Println(girc.StripRaw(out))
	}
}
