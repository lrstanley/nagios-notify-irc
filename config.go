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

const exampleConfig = `# nagios-notify-irc
# generated on: %s
#  NOTE:
#    Most of the entries in this file are NOT required. Simply remove them to
#    have them fallback to their pre-configured defaults.

# The path where the unix socket device file will be stored. This is used by
# both the client and the daemon. Both the client and server must have perms
# to this file, which includes the user which you are invoking the client from
# (likely the "nagios" user).
socket_file = "/tmp/notify-irc.sock"

# The delay in seconds to wait before trying to reconnect after an error.
reconnect_delay = 30

# Default configuration overrides.
default_port = 6667
default_nick = "nagios"
default_name = "Nagios alert relay"
default_user = "nagios"

# Specify as many "[[server]]" blocks as you wish. The only required fields
# are "id" and "hostname", although you should probably fill in at least one
# channel too.

# Below is an example of all of the configuration options.
[[server]]
id = "full-example"
hostname = "irc.example.com"
password = ""
bind = ""
tls = false
tls_skip_verify = false
port = 6667
channels = []
disable_colors = false
nick = "nagios"
name = "Nagios alert relay"
user = "nagios"
sasl_user = ""
sasl_pass = ""

# Below is a shortened but valid example, which also shows how you would
# specify a channel which requires a password to join.
[[server]]
id = "example-1"
hostname = "irc.example2.com"
port = 6697
tls = true
channels = ["#dev", "#secret channel-key"]
`

type GenConfig struct{}

func (GenConfig) Execute(_ []string) error {
	fmt.Printf(exampleConfig+"\n", time.Now().Format(time.ANSIC))
	return nil
}
