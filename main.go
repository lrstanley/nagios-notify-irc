package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	gflags "github.com/jessevdk/go-flags"
)

var version, commit, date = "unknown", "master", "-"

type Flags struct {
	ConfigFile string    `short:"c" long:"config" description:"configuration file location" default:"/etc/notify-irc.toml"`
	Debug      bool      `short:"d" long:"debug" description:"enable debug output"`
	Daemon     Daemon    `command:"daemon" description:"daemon runs and accepts messages for the irc server (generally not run directly)"`
	Client     Client    `command:"client" description:"client connects to a running daemon which forwards messages to the server"`
	GenConfig  GenConfig `command:"gen-config" description:"generate and output an example configuration file"`
}

var flags Flags
var conf *Config
var debug = log.New(ioutil.Discard, "", log.LstdFlags)

func main() {
	parser := gflags.NewParser(&flags, gflags.HelpFlag)
	parser.CommandHandler = func(cmd gflags.Commander, args []string) error {
		if _, ok := cmd.(*GenConfig); ok {
			return cmd.Execute(args)
		}

		var err error
		conf, err = parseConfig(flags.ConfigFile)

		if err != nil {
			return err
		}
		return cmd.Execute(args)
	}
	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if flags.Debug {
		debug.SetOutput(os.Stdout)
	}
}

func catch() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-signals
	fmt.Println("\ninvoked termination, cleaning up")
}
