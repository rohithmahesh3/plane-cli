package integrationtest

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

const minGap = 5 * time.Second

var (
	lockPath  = filepath.Join(os.TempDir(), "plane-cli-integration-rate-limit.lock")
	stampPath = filepath.Join(os.TempDir(), "plane-cli-integration-rate-limit.stamp")
)

// WaitForSlot enforces a minimum delay between integration tests, even across
// separate `go test` package processes, to avoid Plane API rate limits.
func WaitForSlot(t *testing.T) {
	t.Helper()

	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatalf("open integration rate limit lock: %v", err)
	}

	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		_ = lockFile.Close()
		t.Fatalf("lock integration rate limit file: %v", err)
	}

	if data, err := os.ReadFile(stampPath); err == nil {
		lastRunNs, parseErr := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if parseErr == nil {
			wait := minGap - time.Since(time.Unix(0, lastRunNs))
			if wait > 0 {
				t.Logf("integration rate limit guard: sleeping %s", wait.Round(time.Millisecond))
				time.Sleep(wait)
			}
		}
	}

	t.Cleanup(func() {
		now := strconv.FormatInt(time.Now().UnixNano(), 10)
		if err := os.WriteFile(stampPath, []byte(now), 0o600); err != nil {
			t.Fatalf("write integration rate limit stamp: %v", err)
		}
		_ = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
		_ = lockFile.Close()
	})
}
