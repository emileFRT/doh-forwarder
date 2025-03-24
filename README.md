# doh-forwarder

This is a simple but resilient pure-go doh resolver that forward local dns request to Doh services.
Assume this server as work in progress, although usable.
We loosely tried to stick to the [suckless](https://suckless.org) philosophy.

## Requirement

Go tooling is required if building from sources.
The only (pure-go) extra lib is github.com/miekg/dns for handling dns messages.

## How to setup

- (optional) change config.go to your preferences (default should work)
- (optional) change $PREFIX in the makefile
- `make install`
- redirect dns traffic to the server (`nameserver 127.0.0.1` in /etc/resolv.conf on linux for example)
- launch doh-forwarder

You might also want to set up a service manager to control the server such as sv or systemd

### Usage

There is a "-v" cli option that will logs queries, answer and errors.
That option has been imagined to debug and check configurations.
Keep in mind that some programs bypass local dns.

### Non goals:
- caching
- efficient
- secure (we allow failback if not blocked by quad9)
- exhaustive support for all doh services

### Goals:
- simple
- robust
- usable
- use quad9 threat inteligence to block malware
- suckless