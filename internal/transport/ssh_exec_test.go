package transport

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestSSHExecIsRegistered(t *testing.T) {
	tr, err := Get("ssh-exec")
	if err != nil {
		t.Fatalf("ssh-exec transport not registered: %v", err)
	}
	if tr.Name() != "ssh-exec" {
		t.Fatalf("name = %q, want ssh-exec", tr.Name())
	}
}

func TestCleanExecNormalisesCRLFAndAnsi(t *testing.T) {
	in := "\x1b[32m# 2026 by RouterOS\x1b[0m\r\n/interface ether1\r\n"
	got := cleanExec(in)
	want := "# 2026 by RouterOS\n/interface ether1\n"
	if got != want {
		t.Fatalf("cleanExec:\n got %q\nwant %q", got, want)
	}
}

// TestSSHExecLive dials a real device and captures its config over the exec
// transport. Skips unless RUSTED_LIVE_HOST is set, so CI never needs a device.
// Proves the rusted#2 fix end-to-end: MikroTik "/export terse" comes back full.
func TestSSHExecLive(t *testing.T) {
	host := os.Getenv("RUSTED_LIVE_HOST")
	if host == "" {
		t.Skip("set RUSTED_LIVE_HOST/USER/PASS to run against a real device")
	}
	tr, _ := Get("ssh-exec")
	sess, err := tr.Dial(context.Background(), Target{
		Host: host, Username: os.Getenv("RUSTED_LIVE_USER"), Password: os.Getenv("RUSTED_LIVE_PASS"),
	})
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer sess.Close()
	out, err := sess.SendCommand("/export terse")
	if err != nil {
		t.Fatalf("SendCommand: %v", err)
	}
	if len(strings.TrimSpace(out)) == 0 {
		t.Fatal("captured empty configuration (the rusted#2 bug)")
	}
	t.Logf("captured %d bytes; first line: %q", len(out), strings.SplitN(out, "\n", 2)[0])
}
