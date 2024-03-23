package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestCliFunction(t *testing.T) {
	testCases := []struct {
		args    []string
		wantErr bool
	}{
		{[]string{"--version"}, false},
		{[]string{"--help"}, false},
		{[]string{}, true},
		{[]string{"--file", "example/user.go", "--struct", "User"}, false},
		{[]string{"-f", "example/user.go", "-s", "User"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "--force"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "-F"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "--force", "--version"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "-F", "--version"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "--force", "--help"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "-F", "--help"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "--force", "--help", "--version"}, false},
		{[]string{"-f", "example/user.go", "-s", "User", "-F", "--help", "--version"}, false},
	}

	for _, tc := range testCases {
		tempDir, err := os.MkdirTemp("", "testing")
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}
		defer func(path string) {
			err = os.RemoveAll(path)
			if err != nil {
				t.Fatalf("Failed to remove temporary directory: %v", err)
			}
		}(tempDir)

		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		stdout := os.Stdout
		os.Stdout = w

		os.Args = append([]string{"structera"}, tc.args...)

		flagset := flag.NewFlagSet("test", flag.ContinueOnError)
		err = cli(flagset)

		w.Close()
		os.Stdout = stdout
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatal(err)
		}

		if (err != nil) != tc.wantErr {
			fmt.Print(err)
			t.Errorf("cli() with args %v; want error: %v, got error: %v", tc.args, tc.wantErr, err != nil)
		}
	}
}
