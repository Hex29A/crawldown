package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The summary contains the log file path; a path with a percent sign must be
// written verbatim and not be interpreted as a format string.
func TestCloseWritesSummaryVerbatim(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "out%d%s")
	l, err := New(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	l.PagesCrawled = 1
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "crawl.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "out%d%s") {
		t.Fatalf("summary mangled the log path:\n%s", data)
	}
}
