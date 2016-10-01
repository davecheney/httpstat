# httpstat 1.0.0

This project started [a week ago](https://github.com/davecheney/httpstat/commit/29f4578777fdb86c6c0df9a826d047e51bc587f7) as an attempt to replicate the visual presentation of @reorx's [httpstat.py](https://github.com/reorx/httpstat) tool.

From my initial efforts a swam of contributors descended on this project and took it from a proof of concept to a capable tool that is usable across Windows, Linux, and Mac, without any external dependencies.

## That's it folks!

The goal of this project was not to replicate `curl(1)`, but to replicate the visual presentation of `httpstat.py`.
Along the way we've picked up a lot of useful features to round out the general idea of "talk to a server and time the round trip", including being a test bed for the [`httptrace`](https://golang.org/pkg/net/http/httptrace/) package, introduced in Go 1.7.

With the 1.0.0 release, I'm confident that httpstat is a faithful immitation of @reorx's tool, and so I'm declaring this project done.
I'll still be accepting bug reports and will keep this tool up to date with future releses of Go, but no new feature requests will be accepted.

This project is open sourced under a permissive licence, I encourage anyone who wants to hack on it punch that fork button and get coding. Enjoy! 

## Installation
`httpstat` requires Go 1.7.1 or later.
```
% go get -u github.com/davecheney/httpstat
```
## Usage
```
% httpstat
Usage: httpstat [OPTIONS] URL

OPTIONS:
  -E string
        client cert file for tls config
  -H value
        set HTTP header; repeatable: -H 'Accept: ...' -H 'Range: ...'
  -I    don't read body of request
  -L    follow 30x redirects
  -O    save body as remote filename
  -X string
        HTTP method to use (default "GET")
  -d string
        the body of a POST or PUT request
  -k    allow insecure SSL connections
  -o string
        output file for body
  -v    print version number

ENVIRONMENT:
  HTTP_PROXY    proxy for HTTP requests; complete URL or HOST[:PORT]
                used for HTTPS requests if HTTPS_PROXY undefined
  HTTPS_PROXY   proxy for HTTPS requests; complete URL or HOST[:PORT]
  NO_PROXY      comma-separated list of hosts to exclude from proxy
```
## Features

- Windows/BSD/Linux supported.
- HTTP and HTTPS are supported, for self signed certificates use `-k`.
- Skip timing the body of a response with `-I`.
- Follow 30x redirects with `-L`.
- Change HTTP method with `-X METHOD`.
- Provide a `PUT` or `POST` request body with `-d string`. To supply the `PUT` or `POST` body as a file, use `-d @filename`.
- Add extra request headers with `-H 'Name: value'`.
- The response body is usually discarded, you can use `-o filename` to save it to a file, or `-O` to save it to the file name suggested by the server.
- HTTP/HTTPS proxies supported via the usual `HTTP_PROXY`/`HTTPS_PROXY` env vars (as well as lower case variants).
- Supply your own client side certificate with `-E cert.pem`.

## Thanks

This project would not have been possible without the help of testers and contributors who provided feedback, feature requests, bug fixes, documentation fixes, and pull requests. Thank you to:

@amy, @ble, @bogem, @freeformz, @gsquire, @husobee, @imarko, @inkel, @joshi4, @jrozner, @kevinburke, @mattn, @mholt, @mibk, @moorereason, @tcnksm, @theckman, and @Xymist.

I'd like to give special recognition to the contributions of @moorereason who continually sent pull requests and bug fixes for his, and other's features. Thank you.




