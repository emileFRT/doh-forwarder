package main

import "time"

const (
	listenAddr = "127.0.0.1:53"
	timeout    = 5 * time.Second
)

// doh endpoints that support doh wireformat
var dohEndpoints = []string{
	"https://9.9.9.9:5053/dns-query",       // Malware blocking
	"https://149.112.112.9:5053/dns-query", // Malware blocking
	"https://2620:fe::9:5053/dns-query",    // Malware blocking
	"https://9.9.9.9/dns-query",            // Standard security
}

// ips of global dns that will be tried after the systems ones if none are functioning
var fallback_when_default_not_set = []string{
	"9.9.9.9", "149.112.112.112", // Quad9 fallbacks
	"1.1.1.1", // cloudflare
	"8.8.8.8", // google
}
