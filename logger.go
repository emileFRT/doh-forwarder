package main

import (
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	logLevel    int
)

const (
	LogSilent = iota
	LogError
	LogVerbose
)

func init() {
	// Initialize loggers
	infoLogger = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)

	// Set default log level
	logLevel = LogSilent
	if os.Getenv("VERBOSE") != "" {
		logLevel = LogError
	}
}

func setLogLevel(level int) {
	logLevel = level
}

func logQuery(q *dns.Msg, client string) {
	if logLevel < LogVerbose || len(q.Question) == 0 {
		return
	}

	infoLogger.Printf("[QUERY] %s %s %s",
		client,
		dns.TypeToString[q.Question[0].Qtype],
		strings.TrimSuffix(q.Question[0].Name, "."),
	)
}

func logAnswer(a *dns.Msg, client string) {
	if logLevel < LogVerbose {
		return
	}

	if len(a.Answer) > 0 {
		infoLogger.Printf("[ANSWER] %s %d records", client, len(a.Answer))
		for _, rr := range a.Answer {
			infoLogger.Printf("  %s", rr)
		}
	} else {
		infoLogger.Printf("[ANSWER] %s No records", client)
	}
}

func logBlock(q *dns.Msg, client string) {
	if logLevel < LogError {
		return
	}

	qname := "unknown"
	if len(q.Question) > 0 {
		qname = strings.TrimSuffix(q.Question[0].Name, ".")
	}

	errorLogger.Printf("[BLOCK] %s %s",
		client,
		qname,
	)
}

func logError(format string, args ...interface{}) {
	if logLevel >= LogError {
		errorLogger.Printf("[ERROR] "+format, args...)
	}
}
