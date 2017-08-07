package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	ID            string
	Hostname      string
	Password      string
	Bind          string
	Port          int
	TLS           bool
	TLSVerify     bool
	Channels      []string
	DisableColors bool
	Nick          string
	User          string
	Ident         string
	SASLPlain     string
}

type Daemon struct{}

func (s *Daemon) Execute([]string) error {

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("\nexiting")
	return nil
}
