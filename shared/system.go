package shared

import (
	"fmt"
	"net"

	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func SystemUUID() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.HardwareAddr != nil {
			hwAddr := iface.HardwareAddr
			uuid, err := uuid.FromBytes(hwAddr)
			if err != nil {
				return "", err
			}
			return uuid.String(), nil
		}
	}

	return "", fmt.Errorf("no network interface found")
}

func SystemDevice() (string, error) {
	cmd := exec.Command("dmidecode", "-s", "system-product-name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	deviceName := strings.TrimSpace(string(output))
	if deviceName == "" {
		return "", fmt.Errorf("unable to retrieve device name")
	}

	return deviceName, nil
}

func SystemSerial() (string, error) {
	cmd := exec.Command("dmidecode", "-s", "system-serial-number")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	serialNumber := strings.TrimSpace(string(output))
	if serialNumber == "" {
		return "", fmt.Errorf("unable to retrieve serial number")
	}

	return serialNumber, nil
}