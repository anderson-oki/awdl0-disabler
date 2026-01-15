package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anderson-oki/awdl0-disabler/internal/adapters/configuration"
	"github.com/anderson-oki/awdl0-disabler/internal/adapters/filesystem"
	"github.com/anderson-oki/awdl0-disabler/internal/adapters/network"
	"github.com/anderson-oki/awdl0-disabler/internal/adapters/persistence"
	"github.com/anderson-oki/awdl0-disabler/internal/adapters/system"
	"github.com/anderson-oki/awdl0-disabler/internal/adapters/ui"
	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"github.com/anderson-oki/awdl0-disabler/internal/core/services"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	systemAdapter := system.NewSystemAdapter()
	if !systemAdapter.HasElevatedPrivileges() {
		fmt.Println("Error: This application requires elevated privileges (root/sudo) to manage network interfaces.")

		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		os.Exit(1)
	}

	configDirPath := filepath.Join(homeDir, ".config", "awdl0-disabler")
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		fmt.Printf("Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	logsDirPath := filepath.Join(homeDir, ".awdl0-disabler", "logs")
	if err := os.MkdirAll(logsDirPath, 0755); err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	configFilePath := filepath.Join(configDirPath, "config.json")
	configAdapter := configuration.NewJSONConfigAdapter(configFilePath)
	networkAdapter := network.NewShellNetworkAdapter()
	loggerAdapter := filesystem.NewFileLoggerAdapter(logsDirPath)
	repoAdapter := persistence.NewMemoryEventRepo()

	if events, err := loggerAdapter.ReadEvents(time.Now()); err == nil {
		for _, e := range events {
			repoAdapter.Add(e)
		}
	} else {
		fmt.Printf("Warning: Failed to read existing logs: %v\n", err)
	}

	config, err := configAdapter.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)

		os.Exit(1)
	}

	monitorService := services.NewMonitorService(networkAdapter, loggerAdapter, repoAdapter, config)
	statsService := services.NewStatsService(repoAdapter)

	appServices := ui.AppServices{
		Monitor: monitorService,
		Stats:   statsService,
		Config:  config,
		ConfigOnSave: func(c *domain.Config) error {
			return configAdapter.Save(c)
		},
	}

	model := ui.NewModel(appServices)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	_ = monitorService.Restore()
}
