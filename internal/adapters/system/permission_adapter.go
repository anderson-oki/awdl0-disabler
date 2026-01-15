package system

import "os"

type SystemAdapter struct{}

func NewSystemAdapter() *SystemAdapter {
	return &SystemAdapter{}
}

func (a *SystemAdapter) HasElevatedPrivileges() bool {
	return os.Geteuid() == 0
}
