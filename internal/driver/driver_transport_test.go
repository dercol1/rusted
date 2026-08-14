package driver

import "testing"

func TestMikrotikPrefersExecTransport(t *testing.T) {
	d, ok := Get("mikrotik_routeros")
	if !ok {
		t.Fatal("mikrotik_routeros driver not registered")
	}
	if d.Transport != "ssh-exec" {
		t.Fatalf("mikrotik transport = %q, want ssh-exec (rusted#2)", d.Transport)
	}
}
