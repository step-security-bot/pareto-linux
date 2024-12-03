package check

import (
	"fmt"
	"net"
	"time"
)

type Sharing struct {
	passed bool
	ports  map[int]string
}

// Name returns the name of the check
func (f *Sharing) Name() string {
	return "File Sharing is disabled"
}

// Run executes the check
func (f *Sharing) Run() error {
	f.passed = true
	f.ports = make(map[int]string)

	// Samba and NFS ports to check
	shareServices := map[int]string{
		139:  "NetBIOS",
		445:  "SMB",
		2049: "NFS",
		111:  "RPC",
	}

	for port, service := range shareServices {
		addr := fmt.Sprintf("0.0.0.0:%d", port)
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if err == nil {
			f.passed = false
			f.ports[port] = service
			conn.Close()
		}
	}

	return nil
}

// Passed returns the status of the check
func (f *Sharing) Passed() bool {
	return f.passed
}

// CanRun returns whether the check can run
func (f *Sharing) IsRunnable() bool {
	return true
}

// UUID returns the UUID of the check
func (f *Sharing) UUID() string {
	return "b96524e0-850b-4bb8-abc7-517051b6c14e"
}

// ReportIfDisabled returns whether the check should report if it is disabled
func (f *Sharing) ReportIfDisabled() bool {
	return false
}

// Status returns the status of the check
func (f *Sharing) Status() string {
	if !f.Passed() {
		msg := "Sharing services found running on ports:"
		for port, service := range f.ports {
			msg += fmt.Sprintf(" %s(%d)", service, port)
		}
		return msg
	}
	return "No file sharing services found running"
}