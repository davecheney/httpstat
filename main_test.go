package main

import "testing"

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

func TestReadClientCert(t *testing.T) {
	_, err := readClientCert("./test/singlecert.pem")

	if err != nil {
		t.Errorf("unable to read single cert and key: %v", err)
	}

	_, err = readClientCert("./test/multicert.pem")

	if err != nil {
		t.Errorf("unable to read multiple certs and key: %v", err)
	}
}
