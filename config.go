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
	"https://9.9.9.9/dns-query",            // Standard quad9
}
