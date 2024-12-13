package checks

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/log"
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
		8200: "DLNA",
		1900: "Ubuntu Media Sharing",
	}

	for port, service := range shareServices {
		// Check all interfaces
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			// Filter out 127.0.0.1
			if ip.IsLoopback() {
				continue
			}

			address := fmt.Sprintf("%s:%d", ip.String(), port)
			log.WithField("check", f.Name()).WithField("address", address).WithField("protocol", "TCP").Debug("Checking port")
			tcpConn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
			if err == nil {
				defer tcpConn.Close()
				f.passed = false
				f.ports[port] = service
			}
			log.WithField("check", f.Name()).WithField("address", address).WithField("protocol", "UDP").Debug("Checking port")
			udpConn, err := net.DialTimeout("udp", address, 500*time.Millisecond)
			if err == nil {
				defer udpConn.Close()
				f.passed = false
				f.ports[port] = service
			}

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

// PassedMessage returns the message to return if the check passed
func (f *Sharing) PassedMessage() string {
	return "No file sharing services found running"
}

// FailedMessage returns the message to return if the check failed
func (f *Sharing) FailedMessage() string {
	return "Sharing services found running "
}

// RequiresRoot returns whether the check requires root access
func (f *Sharing) RequiresRoot() bool {
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
	return f.PassedMessage()
}