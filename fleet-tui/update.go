package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.leftW = max(28, (msg.Width*58)/100)
		m.headerLines = 1
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.ensureScroll(m.innerH)
			}
		case "down", "j":
			if m.cursor < len(m.vehicles)-1 {
				m.cursor++
				m.ensureScroll(m.innerH)
			}
		case "p":
			m.dispatchCommand("PING", 0)
		case "s":
			m.dispatchCommand("SHUTDOWN", 0)
		case "r":
			m.dispatchCommand("UPDATE_RATE", m.rateMs)
		case "+", "=":
			if m.rateMs < 5000 {
				m.rateMs += 50
			}
		case "-":
			if m.rateMs > 100 {
				m.rateMs -= 50
			}
		}
		return m, nil

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionPress {
			return m, nil
		}
		if msg.Button != tea.MouseButtonLeft {
			return m, nil
		}
		// Header row + top border of the left panel; first vehicle line is two rows below the title bar.
		const mouseContentTop = 2
		row := int(msg.Y)
		col := int(msg.X)
		if row < mouseContentTop || col > m.leftW+4 {
			return m, nil
		}
		idx := row - mouseContentTop + m.scroll
		if idx >= 0 && idx < len(m.vehicles) {
			m.cursor = idx
			m.ensureScroll(m.innerH)
		}
		return m, nil

	case vehiclesDataMsg:
		m.vehicles = msg.data.Vehicles
		if m.cursor >= len(m.vehicles) {
			m.cursor = max(0, len(m.vehicles)-1)
		}
		m.ensureScroll(m.innerH)
		return m, nil

	case metricsDataMsg:
		m.metrics = &msg.data
		return m, nil

	case statusMsg:
		if msg.err != nil {
			m.status = "error: " + msg.err.Error()
		} else {
			m.status = msg.text
		}
		return m, nil
	}

	return m, nil
}

func (m *model) ensureScroll(inner int) {
	if inner <= 0 || len(m.vehicles) == 0 {
		m.scroll = 0
		return
	}

	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}

	if m.cursor >= m.scroll+inner {
		m.scroll = m.cursor - inner + 1
	}

	maxScroll := max(0, len(m.vehicles)-inner)
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
}

func (m *model) selectedID() (uint32, bool) {
	if m.cursor < 0 || m.cursor >= len(m.vehicles) {
		return 0, false
	}
	return m.vehicles[m.cursor].VehicleID, true
}

func (m *model) dispatchCommand(typ string, value uint32) {
	if m.prog == nil {
		return
	}
	id, ok := m.selectedID()
	if !ok {
		m.prog.Send(statusMsg{text: "no vehicle selected"})
		return
	}
	go func() {
		err := postCommand(m.baseURL, id, typ, value)
		if err != nil {
			m.prog.Send(statusMsg{err: err})
			return
		}
		m.prog.Send(statusMsg{text: fmt.Sprintf("%s → #%d", typ, id)})
	}()
}
