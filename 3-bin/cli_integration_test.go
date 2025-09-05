package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

const testTimeout = 60 * time.Second

func projectDir(t *testing.T) string {
	t.Helper()
	// This test file sits in the 3-bin module directory; use working directory from test run
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

// loadDotenv loads KEY from a local .env if present.
func loadDotenv(dir string) {
	dotenvPath := filepath.Join(dir, ".env")
	data, err := os.ReadFile(dotenvPath)
	if err != nil {
		return
	}
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		l := bytes.TrimSpace(line)
		if len(l) == 0 || l[0] == '#' {
			continue
		}
		eq := bytes.IndexByte(l, '=')
		if eq <= 0 {
			continue
		}
		key := string(bytes.TrimSpace(l[:eq]))
		val := string(bytes.TrimSpace(l[eq+1:]))
		if len(val) >= 2 && ((val[0] == '\'' && val[len(val)-1] == '\'') || (val[0] == '"' && val[len(val)-1] == '"')) {
			val = val[1 : len(val)-1]
		}
		if key != "" {
			_ = os.Setenv(key, val)
		}
	}
}

func TestMain(m *testing.M) {
	// Attempt to load .env before running tests
	if dir, err := os.Getwd(); err == nil {
		loadDotenv(dir)
	} else {
		fmt.Fprintln(os.Stderr, "warning: cannot determine working dir for .env loading:", err)
	}
	os.Exit(m.Run())
}

func ensureKeyOrSkip(t *testing.T) {
	t.Helper()
	if os.Getenv("KEY") == "" {
		t.Skip("KEY env is not set; skipping integration test")
	}
}

func runCLI(t *testing.T, args ...string) (stdout string, stderr string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", "."}, args...)...)
	cmd.Dir = projectDir(t)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v\nstdout: %s\nstderr: %s", err, outBuf.String(), errBuf.String())
	}
	return outBuf.String(), errBuf.String()
}

var createdIDRe = regexp.MustCompile(`Created bin id: ([A-Za-z0-9]+)`)

func createBin(t *testing.T) string {
	t.Helper()
	ensureKeyOrSkip(t)
	dataPath := filepath.Join(projectDir(t), "data", "bins.json")
	out, _ := runCLI(t, "-create", "-file", dataPath, "-private")
	m := createdIDRe.FindStringSubmatch(out)
	if len(m) != 2 {
		t.Fatalf("failed to parse created id from output: %q", out)
	}
	return m[1]
}

func deleteBin(t *testing.T, id string) {
	t.Helper()
	ensureKeyOrSkip(t)
	runCLI(t, "-delete", "-id", id)
}

func TestCreateBin(t *testing.T) {
	id := createBin(t)
	// cleanup
	t.Cleanup(func() { deleteBin(t, id) })
}

func TestGetBin(t *testing.T) {
	id := createBin(t)
	t.Cleanup(func() { deleteBin(t, id) })
	out, _ := runCLI(t, "-get", "-id", id)
	if len(out) == 0 {
		t.Fatalf("expected non-empty get output")
	}
}

func TestUpdateBin(t *testing.T) {
	id := createBin(t)
	t.Cleanup(func() { deleteBin(t, id) })
	dataPath := filepath.Join(projectDir(t), "data", "bins.json")
	out, _ := runCLI(t, "-update", "-id", id, "-file", dataPath)
	if out == "" || out[:7] != "Updated" {
		t.Fatalf("expected 'Updated' message, got: %q", out)
	}
}

func TestDeleteBin(t *testing.T) {
	id := createBin(t)
	// delete explicitly; no cleanup needed after successful delete
	out, _ := runCLI(t, "-delete", "-id", id)
	if out == "" || out[:7] != "Deleted" {
		t.Fatalf("expected 'Deleted' message, got: %q", out)
	}
}
