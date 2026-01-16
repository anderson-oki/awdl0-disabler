package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := m.renderHeader()
	footer := m.renderFooter()

	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(footer)
	contentH := m.height - headerH - footerH

	if contentH < 0 {
		contentH = 0
	}

	var content string
	if m.showLogs {
		m.viewport.Width = m.width - m.styles.Logs.GetHorizontalFrameSize()
		m.viewport.Height = contentH - m.styles.Logs.GetVerticalFrameSize()
		content = m.styles.Logs.Render(m.viewport.View())
	} else {
		content = m.renderDashboard(contentH)
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func (m Model) renderHeader() string {
	status := " MONITORING "
	style := m.styles.StatusUp
	if !m.monitoring {
		status = " PAUSED "
		style = m.styles.StatusDown
	}

	awdl0Status := " UNKNOWN "
	awdl0Style := m.styles.StatusUnknown
	if m.awdl0Status == "UP" {
		awdl0Status = " ENABLED "
		awdl0Style = m.styles.StatusUp
	} else if m.awdl0Status == "DOWN" {
		awdl0Status = " DISABLED "
		awdl0Style = m.styles.StatusDown
	}

	content := m.styles.Header.Render(" AWDL0 Disabler ") + " " + style.Render(status) + " " + awdl0Style.Render(awdl0Status) +
		fmt.Sprintf(" Poll: %v (↑/↓: Adjust)", m.services.Config.PollingInterval)

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, content)
}

func (m Model) renderDashboard(availableHeight int) string {
	maxCount := 0
	for _, b := range m.buckets {
		if b.Count > maxCount {
			maxCount = b.Count
		}
	}
	if maxCount == 0 {
		maxCount = 1
	}

	var bars []string

	blocks := []string{" ", " ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

	for _, b := range m.buckets {
		height := 0
		if b.Count <= 0 {
			bars = append(bars, m.styles.Bar.Render(blocks[0]))
			continue
		}

		height = int(float64(b.Count) / float64(maxCount) * 8.0)
		if height > 8 {
			height = 8
		}
		if height < 1 {
			height = 1
		}

		bars = append(bars, m.styles.Bar.Render(blocks[height]))
	}

	graph := strings.Join(bars, "")

	stats := fmt.Sprintf("Last Hour Activity: %d buckets", len(m.buckets))

	dashboardContent := lipgloss.JoinVertical(lipgloss.Center, stats, "\n", graph)

	box := m.styles.Dashboard.Render(dashboardContent)

	if m.awdl0Status == "DOWN" {
		sideEffects := []string{
			"• AirDrop: Disabled",
			"• AirPlay & Sidecar: Impacted",
			"• Continuity & Handoff: Interrupted",
			"• Watch Unlock: Disabled",
		}
		sideEffectsView := m.styles.SideEffects.Render(strings.Join(sideEffects, "  "))
		box = lipgloss.JoinVertical(lipgloss.Center, box, sideEffectsView)
	}

	return lipgloss.Place(m.width, availableHeight, lipgloss.Center, lipgloss.Center, box)
}

func (m Model) renderFooter() string {
	if m.statusMsg != "" {
		style := m.styles.Footer

		if strings.HasPrefix(m.statusMsg, "Error") {
			style = style.Foreground(lipgloss.Color("196")) // Red
		}
		return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, style.Render(m.statusMsg))
	}

	help := "Space: Pause/Resume • L: Logs • E: Toggle awdl0 • ↑: Slower • ↓: Faster • Q: Quit"
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.styles.Footer.Render(help))
}
