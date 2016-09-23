package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	HTTPS_TEMPLATE = `` +
		`  DNS Lookup   TCP Connection   SSL Handshake   Server Processing   Content Transfer` +
		`[   {a0000}  |     {a0001}    |    {a0002}    |      {a0003}      |      {a0004}     ]` +
		`             |                |               |                   |                  |` +
		`    namelookup:{b0000}        |               |                   |                  |` +
		`                        connect:{b0001}       |                   |                  |` +
		`                                    pretransfer:{b0002}           |                  |` +
		`                                                      starttransfer:{b0003}          |` +
		`                                                                                 total:{b0004}`

	HTTP_TEMPLATE = `` +
		`   DNS Lookup   TCP Connection   Server Processing   Content Transfer` + "\n" +
		`[ %7dms  |     %7dms  |        %7dms  |       %7dms  ]` + "\n" +
		`             |                |                   |                  |` + "\n" +
		`    namelookup:%7dms      |                   |                  |` + "\n" +
		`                        connect:%7dms         |                  |` + "\n" +
		`                                      starttransfer:%7dms        |` + "\n" +
		`                                                                 total:%7dms` + "\n"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("usage: %s URL", os.Args[0])
	}

	url, err := url.Parse(args[0])
	if err != nil {
		log.Fatalf("could not parse url %q: %v", args[0], err)
	}

	scheme := url.Scheme
	hostport := url.Host
	host, port := func() (string, int) {
		switch scheme {
		case "https":
			return hostport, 443
		case "http":
			return hostport, 80
		default:
			log.Fatalf("unsupported url scheme %q", scheme)
			return "", 0 // not reached
		}
	}()

	t0 := time.Now() // before dns resolution
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("unable to resolve host: %v", err)
	}

	t1 := time.Now() // after dns resolution, before connect
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalf("unable to connect to host %vv %v", raddr, err)
	}

	t2 := time.Now() // after connect, before request
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Fatalf("unable to create request: %v", err)
	}

	if err := req.Write(conn); err != nil {
		log.Fatalf("failed to write request: %v", err)
	}

	t3 := time.Now() // after request, before read response
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		log.Fatalf("failed to read response: %v", err)
	}

	t4 := time.Now() // after read request, before read body
	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	t5 := time.Now() // after read body
	resp.Body.Close()

	_ = t3

	// print status line and headers
	fmt.Println("\n", resp.Proto, resp.Status)

	for k, v := range resp.Header {
		fmt.Println(k+":", strings.Join(v, ","))
	}

	fmt.Println("\nBody discarded\n")

	switch scheme {
	case "https":
	case "http":
		fmt.Printf(HTTP_TEMPLATE,
			int(t1.Sub(t0)/time.Millisecond), // dns lookup
			int(t2.Sub(t1)/time.Millisecond), // tcp connection
			int(t4.Sub(t2)/time.Millisecond), // server processing
			int(t5.Sub(t4)/time.Millisecond), // content transfer
			int(t1.Sub(t0)/time.Millisecond), // namelookup
			int(t2.Sub(t0)/time.Millisecond), // connect
			int(t4.Sub(t0)/time.Millisecond), // starttransfer
			int(t5.Sub(t0)/time.Millisecond), // total
		)

	}
}
