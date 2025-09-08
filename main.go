package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
)

func main() {
	go ServeDnsUdp()
	ServeDnsTcp()
}

var HTTPClient = &http.Client{}
var logLock = &sync.Mutex{}

func ServeDnsUdp() {
	pck, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer pck.Close()

	buffer := make([]byte, 1500)
	for {
		n, addr, err := pck.ReadFrom(buffer)
		if err != nil {
			logErr(err)
			continue
		}
		go func(data []byte, length int, returnAddr net.Addr) {
			resp := dohProcess(data[:length])
			if resp == nil {
				return
			}
			data, err := io.ReadAll(resp)
			if err != nil {
				logErr(err)
				return
			}
			_, err = pck.WriteTo(data, returnAddr)
			if err != nil {
				logErr(err)
			}
		}(buffer, n, addr)
	}
}

func ServeDnsTcp() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logErr(err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			msg, err := io.ReadAll(conn) // no streams here because several endpoints are allowed so data must be saved in memory anyway
			if err != nil {
				logErr(err)
				return
			}
			resp := dohProcess(msg)
			if resp == nil {
				return
			}
			_, err = io.Copy(conn, resp)
			if err != nil {
				logErr(err)
			}
		}(conn)
	}
}

func dohProcess(msg []byte) io.Reader {
	for _, endpoint := range dohEndpoints {
		req, err := http.NewRequest("POST", endpoint, bytes.NewReader(msg))
		if err != nil {
			logErr(err)
			continue
		}

		// headers for DNS wire format
		req.Header.Set("Content-Type", "application/dns-message")
		req.Header.Set("Accept", "application/dns-message")

		resp, err := HTTPClient.Do(req) // safe to use concurently
		if err != nil {
			logErr(err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			logErr("got http status:", resp.Status, "on endpoint:", endpoint)
			resp.Body.Close()
			continue
		}

		return resp.Body
	}

	logErr("no more endpoint to try, request failure")
	return nil
}

func logErr(msg ...any) {
	logLock.Lock()
	defer logLock.Unlock()

	fmt.Fprintln(os.Stderr, msg...)
}
