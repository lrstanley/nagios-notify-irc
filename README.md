# nagios-notify-irc

Nagios utility for reporting to an IRC channel when an event occurs.

## Table of Contents
- [Releases](#releases)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Daemon](#daemon)
  - [Client](#client)
  - [Colors](#colors)
- [Example Nagios config](#example-nagios-config)
- [License](#license)

## Releases

Check out the [releases](https://github.com/lrstanley/nagios-notify-irc/releases)
page for prebuilt versions. If you need a specific version, feel free to compile
from source (you must install [Go](https://golang.org/doc/install) first):

```
$ git clone https://github.com/lrstanley/nagios-notify-irc.git
$ cd nagios-notify-irc
$ make help
$ make build
```

## Installation

notify-irc should work on Ubuntu, CentOS, and many other distros and
architectures. Below are example commands of how you would install this
(ensure to replace `${VERSION...}` etc, with the appropriate vars):

```
$ wget https://github.com/lrstanley/nagios-notify-irc/releases/download/${VERSION}/nagios-notify-irc_${VERSION_OS_ARCH}.tar.gz
$ tar -C /usr/bin/ -xzvf nagios-notify-irc_${VERSION_OS_ARCH}.tar.gz notify-irc
$ chmod +x /usr/bin/notify-irc
```

## Configuration

Simply run the following command to generate a configuration file which you
can edit.

```
$ notify-irc gen-config > /etc/notify-irc.toml
```

Note that you can add multiple server entries if you wish. You can also place
the configuration file in another location, and specify this location when
you invoke notify-irc like so:

```
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

```
$ notify-irc daemon
```

Though, this should likely be run by systemd or on startup, and not manually
invoked.

### Client

The client is what you will use in your Nagios configurations, which will
pass the message/alert to the daemon. Here is an example:

```
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

```
$ notify-irc client --help
```

### Colors

When passing text to the [client](#client), note that it supports common
irc color codes. You can specify them like `{red}`, `{teal}`, `{bold}`, etc.
To close a color code, you will want to use `{c}`. A full list of supported
codes is shown [here](https://github.com/lrstanley/girc/blob/ef73e5521b5bcbc1248229d8600e574f90a9508d/format.go#L18-L39).

## Example Nagios config

TODO

## License

```
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
```
