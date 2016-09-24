package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	HTTPS_TEMPLATE = `` +
		`  DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer` + "\n" +
		`[%s  |     %s  |    %s  |        %s  |       %s  ]` + "\n" +
		`            |                |               |                   |                  |` + "\n" +
		`   namelookup:%s      |               |                   |                  |` + "\n" +
		`                       connect:%s     |                   |                  |` + "\n" +
		`                                   pretransfer:%s         |                  |` + "\n" +
		`                                                     starttransfer:%s        |` + "\n" +
		`                                                                                total:%s` + "\n"

	HTTP_TEMPLATE = `` +
		`   DNS Lookup   TCP Connection   Server Processing   Content Transfer` + "\n" +
		`[ %s  |     %s  |        %s  |       %s  ]` + "\n" +
		`             |                |                   |                  |` + "\n" +
		`    namelookup:%s      |                   |                  |` + "\n" +
		`                        connect:%s         |                  |` + "\n" +
		`                                      starttransfer:%s        |` + "\n" +
		`                                                                 total:%s` + "\n"
)

var (
	requestBody io.Reader

	grayscale = func(code int) func(string) string {
		if color.NoColor {
			return func(s string) string { return s }
		}
		return func(s string) string {
			return fmt.Sprintf("\x1b[;38;5;%dm%s\x1b[0m", code+232, s)
		}
	}

	// Command line flags.
	httpMethod      string
	postBody        string
	followRedirects bool
	onlyHeader      bool

	usage = fmt.Sprintf("usage: %s URL", os.Args[0])
)

func init() {
	flag.StringVar(&httpMethod, "X", "GET", "HTTP method to use")
	flag.StringVar(&postBody, "d", "", "the body of a POST or PUT request")
	flag.BoolVar(&followRedirects, "L", false, "follow 30x redirects")
	flag.BoolVar(&onlyHeader, "I", false, "don't read body of request")
	flag.Usage = func() {
		os.Stderr.WriteString(usage + "\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf(usage)
	}

	url, err := url.Parse(args[0])
	if err != nil {
		log.Fatalf("could not parse url %q: %v", args[0], err)
	}

	visit(url)
}

// visit visits a url and times the interaction.
// If the response is a 30x, visit follows the redirect.
func visit(url *url.URL) {
	scheme := url.Scheme
	hostport := url.Host
	host, port := func() (string, string) {
		host, port, err := net.SplitHostPort(hostport)
		if err != nil {
			host = hostport
		}
		switch scheme {
		case "https":
			if port == "" {
				port = "443"
			}
		case "http":
			if port == "" {
				port = "80"
			}
		default:
			log.Fatalf("unsupported url scheme %q", scheme)
		}
		return host, port
	}()

	t0 := time.Now() // before dns resolution
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("unable to resolve host: %v", err)
	}

	var conn net.Conn
	t1 := time.Now() // after dns resolution, before connect
	conn, err = net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalf("unable to connect to host %vv %v", raddr, err)
	}

	var t2 time.Time // after connect, before TLS handshake
	if scheme == "https" {
		t2 = time.Now()
		c := tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
		if err := c.Handshake(); err != nil {
			log.Fatalf("unable to negotiate TLS handshake: %v", err)
		}
		conn = c
	}

	t3 := time.Now() // after connect, before request
	if onlyHeader {
		httpMethod = "HEAD"
	}
	if (httpMethod == "POST" || httpMethod == "PUT") && postBody == "" {
		log.Fatal("must supply post body using -d when POST or PUT is used")
	}
	req, err := http.NewRequest(httpMethod, url.String(), strings.NewReader(postBody))
	if err != nil {
		log.Fatalf("unable to create request: %v", err)
	}

	if err := req.Write(conn); err != nil {
		log.Fatalf("failed to write request: %v", err)
	}

	t4 := time.Now() // after request, before read response
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		log.Fatalf("failed to read response: %v", err)
	}
	defer resp.Body.Close()

	t5 := time.Now()
	t6 := t5
	// don't read body if only header is requested
	if !onlyHeader {
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}
		t6 = time.Now() // after read body
	}

	// print status line and headers
	fmt.Printf("\n%s%s%s\n", color.GreenString("HTTP"), grayscale(14)("/"), color.CyanString("%d.%d %s", resp.ProtoMajor, resp.ProtoMinor, resp.Status))

	for k, v := range resp.Header {
		fmt.Println(grayscale(14)(k+":"), color.CyanString(strings.Join(v, ",")))
	}

	fmt.Println("\nBody discarded\n")

	fmta := func(d time.Duration) string {
		return color.CyanString("%7dms", int(d/time.Millisecond))
	}

	fmtb := func(d time.Duration) string {
		return color.CyanString("%-9s", strconv.Itoa(int(d/time.Millisecond))+"ms")
	}

	colorize := func(s string) string {
		v := strings.Split(s, "\n")
		v[0] = grayscale(16)(v[0])
		return strings.Join(v, "\n")
	}

	switch scheme {
	case "https":
		fmt.Printf(colorize(HTTPS_TEMPLATE),
			fmta(t1.Sub(t0)), // dns lookup
			fmta(t2.Sub(t1)), // tcp connection
			fmta(t3.Sub(t2)), // tls handshake
			fmta(t5.Sub(t4)), // server processing
			fmta(t6.Sub(t5)), // content transfer
			fmtb(t1.Sub(t0)), // namelookup
			fmtb(t2.Sub(t0)), // connect
			fmtb(t3.Sub(t0)), // pretransfer
			fmtb(t5.Sub(t0)), // starttransfer
			fmtb(t6.Sub(t0)), // total
		)
	case "http":
		fmt.Printf(colorize(HTTP_TEMPLATE),
			fmta(t1.Sub(t0)), // dns lookup
			fmta(t3.Sub(t1)), // tcp connection
			fmta(t5.Sub(t3)), // server processing
			fmta(t6.Sub(t5)), // content transfer
			fmtb(t1.Sub(t0)), // namelookup
			fmtb(t3.Sub(t0)), // connect
			fmtb(t5.Sub(t0)), // starttransfer
			fmtb(t6.Sub(t0)), // total
		)
	}

	if followRedirects && resp.StatusCode > 299 && resp.StatusCode < 400 {
		loc, err := resp.Location()
		if err != nil {
			if err == http.ErrNoLocation {
				// 30x but no Location to follow, give up.
				return
			}
			log.Fatal("unable to follow redirect: %v", err)
		}
		visit(loc)
	}
}
