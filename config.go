package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	SocketFile string `toml:"socket_file"`

	Servers []Server
}

func parseConfig(path string) (*Config, error) {
	tmpConfig := Config{
		SocketFile: "/tmp/notify-irc.sock",
	}

	_, err := toml.DecodeFile(path, &tmpConfig)
	if err != nil {
		return nil, fmt.Errorf("config[%s]: %s", path, err)
	}

	return &tmpConfig, nil
}

const exampleConfig = `
# The path where the unix socket device file will be stored. This is used by
# both the client and the daemon.
socket_file = "/tmp/notify-irc.sock"
`

type GenConfig struct{}

func (GenConfig) Execute(_ []string) error {
	fmt.Printf("# nagios-notify-irc\n# generated on: %s\n", time.Now().Format(time.ANSIC))
	fmt.Println(exampleConfig)
	return nil
}
