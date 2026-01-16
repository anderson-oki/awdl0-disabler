package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/anderson-oki/awdl0-disabler/internal/core/domain"
	"github.com/anderson-oki/awdl0-disabler/internal/core/services"

	"github.com/charmbracelet/bubbles/viewport"
)

type AppServices struct {
	Monitor      *services.MonitorService
	Stats        *services.StatsService
	Config       *domain.Config
	ConfigOnSave func(*domain.Config) error
}

type Model struct {
	services   AppServices
	monitoring bool
	showLogs   bool
	logBuffer  []domain.Event
	viewport   viewport.Model
	buckets    []domain.Bucket

	// Status Message
	statusMsg string

	// Terminal Dimensions
	width, height int

	// Styles
	styles Styles

	// awdl0 status
	awdl0Status domain.Status
}

func NewModel(services AppServices) Model {
	m := Model{
		services:    services,
		monitoring:  true,
		logBuffer:   []domain.Event{},
		viewport:    viewport.New(0, 0),
		styles:      DefaultStyles(),
		awdl0Status: domain.StatusUnknown,
	}

	// Load historical logs (last 24 hours)
	recentEvents := services.Stats.GetRecentEvents(24 * time.Hour)
	for _, e := range recentEvents {
		m.logBuffer = append(m.logBuffer, e)
	}

	// Trim buffer to max size
	if len(m.logBuffer) > 100 {
		m.logBuffer = m.logBuffer[:100]
	}

	// Initialize viewport content
	m.viewport.SetContent(m.renderLogs())
	m.viewport.GotoBottom()

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.tickCmd(),
	)
}

// Messages
type tickMsg time.Time

type checkResultMsg struct {
	Event *domain.Event
	Err   error
}

type toggleMsg struct {
	Event *domain.Event
	Err   error
}

type configSavedMsg struct {
	Err error
}

type clearStatusMsg struct{}

// Commands
func (m Model) tickCmd() tea.Cmd {
	if !m.monitoring {
		return nil
	}
	// Use the current config interval
	return tea.Tick(m.services.Config.PollingInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) checkNetworkCmd() tea.Cmd {
	return func() tea.Msg {
		event, err := m.services.Monitor.Tick()
		return checkResultMsg{Event: event, Err: err}
	}
}

func (m Model) toggleInterfaceCmd() tea.Cmd {
	return func() tea.Msg {
		event, err := m.services.Monitor.ToggleInterface()

		return toggleMsg{Event: event, Err: err}
	}
}

func (m Model) saveConfigCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.services.ConfigOnSave(m.services.Config)
		return configSavedMsg{Err: err}
	}
}

func clearStatusCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			// Cleanup is handled in main.go after p.Run() returns
			return m, tea.Quit

		case " ":
			m.monitoring = !m.monitoring
			if m.monitoring {
				cmds = append(cmds, m.tickCmd())
			}

		case "l", "L":
			m.showLogs = !m.showLogs

		case "e", "E":
			cmds = append(cmds, m.toggleInterfaceCmd())

		case "up", "right", "down", "left":
			// If logs are shown, these keys are for scrolling the viewport, not changing polling
			if !m.showLogs {
				oldInterval := m.services.Config.PollingInterval
				if msg.String() == "up" || msg.String() == "right" {
					m.services.Config.PollingInterval += 100 * time.Millisecond
				} else { // "down" or "left"
					m.services.Config.PollingInterval -= 100 * time.Millisecond
				}
				m.services.Config.Clamp()
				if oldInterval != m.services.Config.PollingInterval {
					cmds = append(cmds, m.saveConfigCmd())
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 6 // Standardize on dynamic calculation in View()

	case tickMsg:
		if m.monitoring {
			// Trigger next tick
			cmds = append(cmds, m.tickCmd())
			// Trigger network check
			cmds = append(cmds, m.checkNetworkCmd())
		}

	case checkResultMsg:
		// Handle result of network check
		if msg.Event == nil {
			// No event, just update stats
			m.buckets = m.services.Stats.GetHistogram(1*time.Hour, 60)
			return m, nil
		}

		// Event occurred
		m.logBuffer = append(m.logBuffer, *msg.Event)
		// Trim buffer
		if len(m.logBuffer) > 100 {
			m.logBuffer = m.logBuffer[len(m.logBuffer)-100:]
		}
		m.viewport.SetContent(m.renderLogs())
		m.viewport.GotoBottom()

		// Update stats after check
		m.buckets = m.services.Stats.GetHistogram(1*time.Hour, 60)

		// Update awdl0 status
		if msg.Event.Type == domain.EventDisable {
			m.awdl0Status = domain.StatusDown
		} else if msg.Event.Type == domain.EventEnable {
			m.awdl0Status = domain.StatusUp
		}

	case toggleMsg:
		if msg.Err != nil {
			m.statusMsg = fmt.Sprintf("Error toggling: %v", msg.Err)
			cmds = append(cmds, clearStatusCmd())
			return m, tea.Batch(cmds...)
		}

		if msg.Event != nil {
			m.logBuffer = append(m.logBuffer, *msg.Event)
			if len(m.logBuffer) > 100 {
				m.logBuffer = m.logBuffer[len(m.logBuffer)-100:]
			}
			m.viewport.SetContent(m.renderLogs())
			m.viewport.GotoBottom()
			if msg.Event.Type == domain.EventDisable {
				m.awdl0Status = domain.StatusDown
			} else if msg.Event.Type == domain.EventEnable {
				m.awdl0Status = domain.StatusUp
			}
		}

		// Update stats after check
		m.buckets = m.services.Stats.GetHistogram(1*time.Hour, 60)

	case configSavedMsg:
		if msg.Err != nil {
			m.statusMsg = fmt.Sprintf("Error saving: %v", msg.Err)
		} else {
			m.statusMsg = "Settings Saved"
		}
		cmds = append(cmds, clearStatusCmd())

	case clearStatusMsg:
		m.statusMsg = ""
	}

	// Handle viewport updates if logs are shown
	if m.showLogs {
		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) renderLogs() string {
	var content strings.Builder
	for _, event := range m.logBuffer {
		timestamp := m.styles.Timestamp.Render(event.Timestamp.Format("15:04:05"))
		content.WriteString(fmt.Sprintf("%s %s\n", timestamp, event.Message))
	}
	return content.String()
}
