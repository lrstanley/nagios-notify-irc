package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Client struct{}

func (c *Client) Execute([]string) error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("\nexiting")
	return nil
}
