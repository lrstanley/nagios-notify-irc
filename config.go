package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	SocketFile     string   `toml:"socket_file"`
	ReconnectDelay int      `toml:"reconnect_delay"`
	DefaultPort    int      `toml:"default_port"`
	DefaultNick    string   `toml:"default_nick"`
	DefaultName    string   `toml:"default_name"`
	DefaultUser    string   `toml:"default_user"`
	Servers        []Server `toml:"server"`
}

func parseConfig(path string) (*Config, error) {
	tmpConfig := Config{
		SocketFile:     "/tmp/notify-irc.sock",
		ReconnectDelay: 45,
		DefaultPort:    6667,
		DefaultNick:    "nagios",
		DefaultName:    "Nagios alert relay",
		DefaultUser:    "nagios",
	}

	_, err := toml.DecodeFile(path, &tmpConfig)
	if err != nil {
		return nil, fmt.Errorf("config[%s]: %s", path, err)
	}

	if tmpConfig.ReconnectDelay < 10 {
		tmpConfig.ReconnectDelay = 10
	}

	if len(tmpConfig.Servers) == 0 {
		return nil, fmt.Errorf("config[%s]: no servers specified", path)
	}

	servers := []string{}
	for i := 0; i < len(tmpConfig.Servers); i++ {
		for j := 0; j < len(servers); j++ {
			if tmpConfig.Servers[i].ID == servers[j] {
				return nil, fmt.Errorf("config[%s]: duplicate server ID found: %s", path, tmpConfig.Servers[i].ID)
			}
		}
		servers = append(servers, tmpConfig.Servers[i].ID)
	}

	return &tmpConfig, nil
}

const exampleConfig = `
# The path where the unix socket device file will be stored. This is used by
# both the client and the daemon.
socket_file = "/tmp/notify-irc.sock"

reconnect_delay = "30s"

default_port = 6667
default_nick = "nagios"
default_name = "Nagios alert relay"
default_user = "nagios"

[[server]]
id = "test1"
nick = "nagios1"
hostname = "irc.byteirc.org"
port = 6697
tls = true
channels = ["#dev", "#dev1"]

[[server]]
id = "test2"
nick = "nagios2"
hostname = "irc.byteirc.org"
port = 6697
tls = true
channels = ["#dev yourkey", "#dev1"]
`

type GenConfig struct{}

func (GenConfig) Execute(_ []string) error {
	fmt.Printf("# nagios-notify-irc\n# generated on: %s\n", time.Now().Format(time.ANSIC))
	fmt.Println(exampleConfig)
	return nil
}
