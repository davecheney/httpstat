# httpstat [![Build Status](https://github.com/davecheney/httpstat/actions/workflows/push.yml/badge.svg)](https://github.com/davecheney/httpstat/actions/workflows/push.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/davecheney/httpstat)](https://goreportcard.com/report/github.com/davecheney/httpstat)

![Shameless](./screenshot.png)

Imitation is the sincerest form of flattery.

But seriously, https://github.com/reorx/httpstat is the new hotness, and this is a shameless rip off.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
- [Examples](#examples)
- [Use Cases](#use-cases)
- [Running Tests](#running-tests)
- [Building the Project](#building-the-project)
- [CI/CD with GitHub Actions](#cicd-with-github-actions)
- [Contributing](#contributing)

## Installation
`httpstat` requires Go 1.20 or later.
```
go install github.com/davecheney/httpstat@latest
```

### Detailed Installation Instructions

#### Windows
1. Download and install Go from the official website: https://golang.org/dl/
2. Open Command Prompt and run the following command:
   ```
   go install github.com/davecheney/httpstat@latest
   ```

#### macOS
1. Download and install Go from the official website: https://golang.org/dl/
2. Open Terminal and run the following command:
   ```
   go install github.com/davecheney/httpstat@latest
   ```

#### Linux
1. Download and install Go from the official website: https://golang.org/dl/
2. Open Terminal and run the following command:
   ```
   go install github.com/davecheney/httpstat@latest
   ```

### Troubleshooting Installation Issues
- Ensure that Go is installed correctly by running `go version` in your terminal or command prompt.
- Make sure your `GOPATH` and `GOROOT` environment variables are set correctly.
- If you encounter any issues, refer to the official Go installation guide: https://golang.org/doc/install

## Usage
```
httpstat https://example.com/
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

## Examples

### Basic Usage
```
httpstat https://example.com/
```

### Using Custom HTTP Method
```
httpstat -X POST https://example.com/
```

### Providing a Request Body
```
httpstat -X POST -d "key=value" https://example.com/
```

### Adding Extra Request Headers
```
httpstat -H "Authorization: Bearer <token>" https://example.com/
```

## Use Cases

### Measuring Response Time
Use `httpstat` to measure the response time of your web application and identify performance bottlenecks.

### Debugging HTTP Requests
`httpstat` can help you debug HTTP requests by providing detailed timing information for each phase of the request.

### Monitoring API Performance
Monitor the performance of your APIs by regularly running `httpstat` and tracking the response times.

## Running Tests
To run tests, use the following command:
```
go test -race ./...
```

## Building the Project
To build the project using the provided `Makefile`, use the following command:
```
make
```

## CI/CD with GitHub Actions
This project uses GitHub Actions for CI/CD. The workflow is defined in the `.github/workflows/push.yml` file. It runs tests, performs static analysis, and verifies the Go modules.

## Contributing

Bug reports are most welcome, but with the exception of #5, this project is closed.

Pull requests must include a `fixes #NNN` or `updates #NNN` comment. 

Please discuss your design on the accompanying issue before submitting a pull request. If there is no suitable issue, please open one to discuss the feature before slinging code. Thank you.
