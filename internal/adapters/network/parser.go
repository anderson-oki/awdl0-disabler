package network

import (
	"strings"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
)

// ParseInterfaceStatus analyzes ifconfig output to determine if interface is UP
func ParseInterfaceStatus(output string) domain.Status {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if !strings.Contains(line, "flags=") {
			continue
		}

		start := strings.Index(line, "<")
		end := strings.Index(line, ">")

		if start == -1 || end == -1 || start >= end {
			continue
		}

		flags := line[start+1 : end]
		if strings.Contains(flags, "UP") {
			return domain.StatusUp
		}
	}
	return domain.StatusDown
}
