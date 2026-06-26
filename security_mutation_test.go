package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestProviderHeaderLimitMutations(t *testing.T) {
	mutations := []struct {
		name        string
		path        string
		original    string
		replacement string
	}{
		{"header limit", "main.go", "const maxTwilioResponseHeaderBytes = 64 * 1024", "const maxTwilioResponseHeaderBytes = 0"},
		{"transport assignment", "main.go", "transport.MaxResponseHeaderBytes = maxTwilioResponseHeaderBytes", "transport.MaxResponseHeaderBytes = 0"},
		{"real send routing", "main.go", "newTwilioRestClient(config, newTwilioTransport())", "newTwilioRestClient(config, http.DefaultTransport)"},
		{"regression test", "main_test.go", "TestNewTwilioTransportBoundsResponseHeadersWithoutMutatingDefault", "RemovedTwilioTransportHeaderLimitRegression"},
		{"completed plan", "docs/plans/2026-06-26-twilio-response-header-limit.md", "status: completed", "status: in progress"},
	}

	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			fixture := t.TempDir()
			copyMutationFixture(t, fixture)
			path := filepath.Join(fixture, mutation.path)
			contents, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Count(string(contents), mutation.original) != 1 {
				t.Fatalf("mutation anchor drifted: %s", mutation.name)
			}
			updated := strings.Replace(string(contents), mutation.original, mutation.replacement, 1)
			if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
				t.Fatal(err)
			}

			command := exec.Command("/bin/sh", "scripts/check-baseline.sh")
			command.Dir = fixture
			if err := command.Run(); err == nil {
				t.Fatalf("hostile mutation survived: %s", mutation.name)
			}
		})
	}
}

func copyMutationFixture(t *testing.T, destination string) {
	t.Helper()
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == ".git" || path == ".explore" {
			return filepath.SkipDir
		}
		if path == "." {
			return nil
		}
		target := filepath.Join(destination, path)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, contents, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
