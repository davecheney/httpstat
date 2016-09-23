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
	"regexp"
	"runtime"
	"strconv"
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
		`[ %s  |     %s  |        %s  |       %s  ]` + "\n" +
		`             |                |                   |                  |` + "\n" +
		`    namelookup:%s      |                   |                  |` + "\n" +
		`                        connect:%s         |                  |` + "\n" +
		`                                      starttransfer:%s        |` + "\n" +
		`                                                                 total:%s` + "\n"
)

const ISATTY = runtime.GOOS != "windows" // TODO(dfc) make this respect the actual pty state

func makeColor(code int) func(string) string {
	if !ISATTY {
		return func(s string) string { return s }
	}
	return func(s string) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, s)
	}
}

var (
	red       = makeColor(31)
	green     = makeColor(32)
	yellow    = makeColor(33)
	blue      = makeColor(34)
	magenta   = makeColor(35)
	cyan      = makeColor(36)
	bold      = makeColor(1)
	underline = makeColor(4)

	grayscale = func(code int) func(string) string {
		if !ISATTY {
			return func(s string) string { return s }
		}
		return func(s string) string {
			return fmt.Sprintf("\x1b[;38;5;%dm%s\x1b[0m", code+232, s)
		}
	}
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

	if portInHost := regexp.MustCompile(":\\d+").FindString(host); portInHost != "" {
		portNum := regexp.MustCompile(":").Split(portInHost, -1)[0]

		port, err = strconv.Atoi(portNum)

		if err != nil {
			log.Fatalf("unable to parse port in host: %v", err)
		}
	}

	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("unable to resolve host: %v", err)
	}

	t1 := time.Now() // after dns resolution, before connect
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalf("unable to connect to host %v %v", raddr, err)
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
	fmt.Printf("\n%s%s%s\n", green("HTTP"), grayscale(14)("/"), cyan(fmt.Sprintf("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, resp.Status)))

	for k, v := range resp.Header {
		fmt.Println(grayscale(14)(k+":"), cyan(strings.Join(v, ",")))
	}

	fmt.Println("\nBody discarded\n")

	fmta := func(d time.Duration) string {
		return cyan(fmt.Sprintf("%7dms", int(d/time.Millisecond)))
	}

	fmtb := func(d time.Duration) string {
		return cyan(fmt.Sprintf("%-9s", strconv.Itoa(int(d/time.Millisecond))+"ms"))
	}

	colorize := func(s string) string {
		v := strings.Split(s, "\n")
		v[0] = grayscale(16)(v[0])
		return strings.Join(v, "\n")
	}

	switch scheme {
	case "https":
		// TODO(dfc) handle HTTPS
	case "http":
		fmt.Printf(colorize(HTTP_TEMPLATE),
			fmta(t1.Sub(t0)), // dns lookup
			fmta(t2.Sub(t1)), // tcp connection
			fmta(t4.Sub(t2)), // server processing
			fmta(t5.Sub(t4)), // content transfer
			fmtb(t1.Sub(t0)), // namelookup
			fmtb(t2.Sub(t0)), // connect
			fmtb(t4.Sub(t0)), // starttransfer
			fmtb(t5.Sub(t0)), // total
		)

	}
}
