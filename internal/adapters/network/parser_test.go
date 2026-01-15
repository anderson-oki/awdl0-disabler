package network_test

import (
	"github.com/anderson-oki/awdl0-disabler/internal/adapters/network"
	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"testing"
)

func TestParseInterfaceStatus(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected domain.Status
	}{
		{
			name:     "Interface UP",
			output:   "awdl0: flags=8051<UP,POINTOPOINT,RUNNING,MULTICAST> mtu 1500\n\tether ...",
			expected: domain.StatusUp,
		},
		{
			name:     "Interface DOWN",
			output:   "awdl0: flags=8050<POINTOPOINT,RUNNING,MULTICAST> mtu 1500\n\tether ...",
			expected: domain.StatusDown,
		},
		{
			name:     "Empty Output",
			output:   "",
			expected: domain.StatusDown,
		},
		{
			name:     "Complex Output",
			output:   "lo0: flags=8049<UP,LOOPBACK,RUNNING,MULTICAST> mtu 16384\nawdl0: flags=8051<UP,POINTOPOINT,RUNNING,MULTICAST> mtu 1500",
			expected: domain.StatusUp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := network.ParseInterfaceStatus(tt.output)
			if got != tt.expected {
				t.Errorf("ParseInterfaceStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}
