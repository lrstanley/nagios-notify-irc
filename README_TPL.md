# nagios-notify-irc

Nagios utility for reporting to an IRC channel when an event occurs.

_A sister project to [nagios-check-ircd](https://github.com/lrstanley/nagios-check-ircd)_

## Table of Contents
- [Installation](#installation)
  - [Ubuntu/Debian](#ubuntudebian)
  - [CentOS/Redhat](#centosredhat)
  - [Manual Install](#manual-install)
  - [Source](#source)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Daemon](#daemon)
  - [Client](#client)
    - [Templating](#templating)
    - [Colors](#colors)
- [Example Nagios config](#example-nagios-config)
- [Contributing](#contributing)
- [License](#license)

## Installation

Check out the [releases](https://github.com/lrstanley/nagios-notify-irc/releases)
page for prebuilt versions. notify-irc should work on ubuntu/debian,
centos/redhat/fedora, etc. Below are example commands of how you would install
the utility.

**NOTE**: If you are running nagios as a different user, you _will_ need to
update the service files to the correct user.

### Ubuntu/Debian

```console
$ wget https://liam.sh/ghr/notify-irc_[[tag]]_[[os]]_[[arch]].deb
$ dpkg -i notify-irc_[[tag]]_[[os]]_[[arch]].deb
$ notify-irc gen-config > /etc/notify-irc.toml # may want to edit the config as well
$ systemctl enable notify-irc
$ systemctl start notify-irc
```

### CentOS/Redhat

```console
$ yum localinstall https://liam.sh/ghr/notify-irc_[[tag]]_[[os]]_[[arch]].rpm
$ notify-irc gen-config > /etc/notify-irc.toml # may want to edit the config as well
$ systemctl enable notify-irc
$ systemctl start notify-irc
```

### Manual Install

```console
$ wget https://liam.sh/ghr/notify-irc_[[tag]]_[[os]]_[[arch]].tar.gz
$ tar -C /usr/bin/ -xzvf notify-irc_[[tag]]_[[os]]_[[arch]].tar.gz notify-irc
$ chmod +x /usr/bin/notify-irc
$ notify-irc gen-config > /etc/notify-irc.toml # may want to edit the config as well
$ notify-irc daemon # run this in a screen, cron, your own init script, etc.
```

### Source

If you need a specific version, feel free to compile from source (you must
install [Go](https://golang.org/doc/install) first):

```console
$ git clone https://github.com/lrstanley/nagios-notify-irc.git
$ cd nagios-notify-irc
$ make help
$ make build
```

## Configuration

Simply run the following command to generate a configuration file which you
can edit.

```console
$ notify-irc gen-config > /etc/notify-irc.toml
```

Note that you can add multiple server entries if you wish. You can also place
the configuration file in another location, and specify this location when
you invoke notify-irc like so:

```console
$ notify-irc -c path/to/your/config.toml [FLAGS] [SUB-COMMAND] [ARGS]
```

## Usage

The way notify-irc works, is that it runs a daemon in the background which
stays connected to irc. When an alert comes in, the notify-irc client connects
to the daemon, and forwards the message to the necessary server/channel.

This is done to prevent unwanted join/part spam from the bot, and give a good
representation of knowing that the alert bot is still functioning.

### Daemon

To run the daemon, simply execute the following:

```console
$ notify-irc daemon
```

Though, this should likely be run by systemd or on startup, and not manually
invoked.

### Client

The client is what you will use in your Nagios configurations, which will
pass the message/alert to the daemon. Here is an example:

```console
$ notify-irc -s example-1 -p "@" -c "#your-channel" -c "#another-channel" "Example message" "More info"
```

With the above command, this will send an alert to the example-1 server
(specified as the "id" in the configuration file), to the two specified
channels, ping the ops in each of those channels, and send two messages
separate messages to the channel, `Example message` and `More info`.

Note that specifying the server id, ping list, etc, is all optional. This
allows you to generate a message which is sent to multiple networks, across
multiple channels.

See the following for more information:

```console
$ notify-irc client --help
```

#### Templating

Text passed into the [client](#client) by default will allow Go-based
`text/template` templates. For example:

```console
$ notify-irc client -c "#your-channel" '{{ if eq "$SERVICESTATUS" "OK" }}Healthy!{{ else }}Uhoh!{{ end }}: Other stuff here.'
```

**Note**: This can be disabled by passing `--no-tmpl` to `notify-irc client`,
which will pass all input text as plaintext.

#### Colors

When passing text to the [client](#client), note that it supports common
irc color codes. You can specify them like `{red}`, `{teal}`, `{bold}` (or
`{b}` for short), etc. To close a color code, you will want to use `{c}` (
`{red}ERROR{c}: default color with some {b}bold!{b}`). A full list of
supported codes is shown [here](https://github.com/lrstanley/girc/blob/ef73e5521b5bcbc1248229d8600e574f90a9508d/format.go#L18-L39).

## Example Nagios Config

```conf
define command {
	command_name notify_irc_service
	command_line /usr/local/bin/notify-irc client -p '@' -c '*' '[{{if eq "$SERVICESTATE$" "OK"}}{green}{{else}}{red}{{end}}{b}$SERVICESTATE${b}{c}] {yellow}{b}$SERVICEDESC${b}{c} :: {teal}$HOSTNAME${c} ({teal}$HOSTADDRESS${c}) :: ({b}$SERVICESTATETYPE${b}: for {cyan}$SERVICEDURATION${c})' '$SERVICEOUTPUT$'
}

define command {
	command_name notify_irc_host
	command_line /usr/local/bin/notify-irc client -p '@' -c '*' '[{{if eq "$HOSTSTATE$" "OK"}}{green}{{else}}{red}{{end}}{b}$HOSTSTATE${b}{c}] {teal}$HOSTNAME${c} ({teal}$HOSTADDRESS${c}) :: ({b}$HOSTSTATETYPE${b}: for {cyan}$HOSTDURATION${c}) :: [ {green}{b}OK:{b} $TOTALHOSTSERVICESOK${c} | {yellow}{b}WARN:{b} $TOTALHOSTSERVICESWARNING${c} | {b}UNKN:{b} $TOTALHOSTSERVICESUNKNOWN$ | {red}{b}CRIT:{b} $TOTALHOSTSERVICESCRITICAL${c} ]' '$HOSTOUTPUT$'
}
```

## Contributing

Please review the [CONTRIBUTING](https://github.com/lrstanley/nagios-notify-irc/blob/master/CONTRIBUTING.md)
doc for submitting issues/a guide on submitting pull requests and helping out.

## License

    LICENSE: The MIT License (MIT)
    Copyright (c) 2017 Liam Stanley <me@liamstanley.io>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
