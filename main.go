package main

import (
	"bytes"
	"flag"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/miekg/dns"
)

var verbose bool

func main() {
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	if verbose {
		setLogLevel(LogVerbose)
	}

	startServer("udp")
	startServer("tcp")
}

func startServer(network string) {
	server := &dns.Server{
		Addr:    listenAddr,
		Net:     network,
		Handler: handler{},
	}
	if err := server.ListenAndServe(); err != nil {
		logError("Server failed: %v", err)
		os.Exit(1)
	}
}

type handler struct{}

func (h handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	client := w.RemoteAddr().String()
	logQuery(r, client)

	resp, blocked := tryDoh(r)
	if blocked {
		logBlock(resp, client)
		w.WriteMsg(resp)
		return
	}

	if resp != nil {
		logAnswer(resp, client)
		w.WriteMsg(resp)
		return
	}

	// Fallback handling
	if resp = tryFallback(r); resp != nil {
		logAnswer(resp, client)
		w.WriteMsg(resp)
		return
	}

	w.WriteMsg(errorResponse(r))
}

func tryDoh(q *dns.Msg) (*dns.Msg, bool) {
	buf, _ := q.Pack()

	for _, endpoint := range dohEndpoints {
		req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/dns-message")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}

		buf, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		answer := new(dns.Msg)
		if err := answer.Unpack(buf); err != nil {
			continue
		}

		// Return early for blocks without verbose parsing
		if answer.Rcode == dns.RcodeNameError && len(answer.Ns) == 0 {
			return answer.SetReply(q), true
		}
		return answer.SetReply(q), false
	}
	return nil, false
}

func tryFallback(q *dns.Msg) *dns.Msg {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil || len(conf.Servers) == 0 {
		return nil
	}

	// Filter out localhost while preserving order
	filtered := make([]string, 0, len(conf.Servers))
	for _, srv := range conf.Servers {
		if srv != "127.0.0.1" && srv != "::1" {
			filtered = append(filtered, srv)
		}
	}

	// Use hardcoded fallbacks if all servers were localhost
	if len(filtered) == 0 {
		filtered = fallback_when_default_not_set
	}

	c := new(dns.Client)
	for _, srv := range filtered {
		serverAddr := net.JoinHostPort(srv, conf.Port)
		resp, _, err := c.Exchange(q, serverAddr)
		if err == nil {
			return resp.SetReply(q)
		}
	}
	return nil
}

func errorResponse(q *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetRcode(q, dns.RcodeServerFailure)
	return m
}
