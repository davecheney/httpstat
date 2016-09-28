package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"https://golang.org", "https://golang.org"},
		{"https://golang.org:443/test", "https://golang.org:443/test"},
		{"localhost:8080/test", "https://localhost:8080/test"},
		{"localhost:80/test", "http://localhost:80/test"},
		{"//localhost:8080/test", "https://localhost:8080/test"},
		{"//localhost:80/test", "http://localhost:80/test"},
	}

	for _, test := range tests {
		u := parseURL(test.in)
		if u.String() != test.want {
			t.Errorf("Given: %s\nwant: %s\ngot: %s", test.in, test.want, u.String())
		}
	}
}

func fmtArgsFail(args []string, name string, want, got interface{}) string {
	return fmt.Sprintf("Given arguments %#v the %#v field was expected to be %#v but was %#v", args, name, want, got)
}

func Test_binArgs_parse(t *testing.T) {
	const (
		//arg0 is used by the runtime
		Arg0    = "/usr/local/bin/httpstat"
		testURL = "https://httpstat.test"
	)

	var out string
	var err error

	cli := []string{Arg0}

	tests := []struct {
		args        []string
		shouldErr   bool
		expectedOut string

		httpMethod, postBody, outputFile, urlString                string
		followRedirects, onlyHeader, insecure, saveOutput, version bool
		httpHeaders                                                headers
	}{
		// bare usage
		{args: cli, shouldErr: true},

		// usage with a URL, make sure insecure is false by default!
		{args: append(cli, testURL), urlString: "https://httpstat.test", insecure: false},

		// -d / --data
		{args: append(cli, "-d", "testData", testURL), postBody: "testData"},
		{args: append(cli, "--data", "testData", testURL), postBody: "testData"},

		// -H / --header
		{
			args:        append(cli, "-H", "Host: httpstat.test", "-H", "Accept: passing-tests", testURL),
			httpHeaders: []string{"Host: httpstat.test", "Accept: passing-tests"},
		},
		{
			args:        append(cli, "--header", "Host: httpstat.test", "--header", "Accept: passing-tests", testURL),
			httpHeaders: []string{"Host: httpstat.test", "Accept: passing-tests"},
		},

		// -I / --head
		{args: append(cli, "-I", testURL), onlyHeader: true},
		{args: append(cli, "--head", testURL), onlyHeader: true},

		// -k / --insecure
		{args: append(cli, "-k", testURL), insecure: true},
		{args: append(cli, "--insecure", testURL), insecure: true},

		// -L / --location
		{args: append(cli, "-L", testURL), followRedirects: true},
		{args: append(cli, "--location", testURL), followRedirects: true},

		// -o / --output
		{args: append(cli, "-o", "/dev/null", testURL), outputFile: "/dev/null"},
		{args: append(cli, "--output", "/dev/null", testURL), outputFile: "/dev/null"},

		// -O / --remote-file
		{args: append(cli, "-O", testURL), saveOutput: true},
		{args: append(cli, "--remote-name", testURL), saveOutput: true},

		// -V / --version
		{
			args:        append(cli, "-V"),
			expectedOut: fmt.Sprintf("/usr/local/bin/httpstat %s (runtime: %s)\n", version, runtime.Version()),
		},
		{
			args:        append(cli, "--version"),
			expectedOut: fmt.Sprintf("/usr/local/bin/httpstat %s (runtime: %s)\n", version, runtime.Version()),
		},

		// -X / --request
		{args: append(cli, "-X", "POST", testURL), httpMethod: "POST"},
		{args: append(cli, "--request", "POST", testURL), httpMethod: "POST"},
	}

	//
	// Massive test runner
	//
	for _, test := range tests {
		args := &binArgs{}

		out, err = args.parse(test.args)

		// should that have errored?
		if test.shouldErr && err == nil {
			t.Errorf("Given arguments %#v an error was expected", test.args)
		}

		// was that an erroneous error?
		if !test.shouldErr && err != nil {
			t.Errorf("Given arguments %#v an error was not expected, got %#v", test.args, err)
		}

		// we were only checking for an error
		// don't continue
		if test.shouldErr {
			continue
		}

		if out != test.expectedOut {
			t.Errorf("Given arguments %#v the following output was expected %#v got %#v", test.args, test.expectedOut, out)
		}

		// we were only checking for help/version output
		// don't continue
		if len(test.expectedOut) > 0 {
			continue
		}

		if test.httpMethod == "" {
			test.httpMethod = "GET"
		}

		if test.urlString == "" {
			test.urlString = testURL
		}

		if aURL := args.URL.String(); aURL != test.urlString {
			t.Errorf("Given arguments %#v the URL was expected to be %#v but was %#v", test.args, test.urlString, aURL)
		}

		//
		// Test the HTTPHeaders field
		//
		if len(args.HTTPHeaders) != len(test.httpHeaders) {
			t.Errorf("Given arguments %#v there was a mismatch in the length of HTTP headers. Expected: %#v; Got: %#v", test.args, test.httpHeaders, args.HTTPHeaders)
		}

		for i, hdr := range test.httpHeaders {
			if args.HTTPHeaders[i] != hdr {
				t.Errorf("Given arguments %#v the header at index %d did not match. Expected %#v Got: %#v", test.args, i, hdr, args.HTTPHeaders[i])
			}
		}

		//
		// Test the remaining fields
		//
		if args.HTTPMethod != test.httpMethod {
			t.Errorf(fmtArgsFail(test.args, "HTTPMethod", test.httpMethod, args.HTTPMethod))
		}
		if args.PostBody != test.postBody {
			t.Errorf(fmtArgsFail(test.args, "PostBody", test.postBody, args.PostBody))
		}
		if args.OutputFile != test.outputFile {
			t.Errorf(fmtArgsFail(test.args, "OutputFile", test.outputFile, args.OutputFile))
		}
		if args.FollowRedirects != test.followRedirects {
			t.Errorf(fmtArgsFail(test.args, "FollowRedirects", test.followRedirects, args.FollowRedirects))
		}
		if args.OnlyHeader != test.onlyHeader {
			t.Errorf(fmtArgsFail(test.args, "OnlyHeader", test.onlyHeader, args.OnlyHeader))
		}
		if args.Insecure != test.insecure {
			t.Errorf(fmtArgsFail(test.args, "Insecure", test.insecure, args.Insecure))
		}
		if args.SaveOutput != test.saveOutput {
			t.Errorf(fmtArgsFail(test.args, "SaveOutput", test.saveOutput, args.SaveOutput))
		}
		if args.Version != test.version {
			t.Errorf(fmtArgsFail(test.args, "Version", test.version, args.Version))
		}
	}
}
