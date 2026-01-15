package network

import (
	"os/exec"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

type ShellNetworkAdapter struct{}

func NewShellNetworkAdapter() *ShellNetworkAdapter {
	return &ShellNetworkAdapter{}
}

func (a *ShellNetworkAdapter) CheckInterface(name string) (domain.Status, error) {
	cmd := exec.Command("ifconfig", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return domain.StatusDown, err
	}

	return ParseInterfaceStatus(string(output)), nil
}

func (a *ShellNetworkAdapter) DisableInterface(name string) error {
	cmd := exec.Command("sudo", "ifconfig", name, "down")
	return cmd.Run()
}

func (a *ShellNetworkAdapter) EnableInterface(name string) error {
	cmd := exec.Command("sudo", "ifconfig", name, "up")
	return cmd.Run()
}
