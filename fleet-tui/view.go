package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevinrestrepoh/fleet-tui/util"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorSubtle   = lipgloss.Color("#6B7280") // muted gray
	colorSelected = lipgloss.Color("#E5E7EB") // soft white (NOT pure white)
	colorOffline  = lipgloss.Color("#9CA3AF") // dim gray
	colorMetricK  = lipgloss.Color("#9CA3AF")
	colorBorder   = lipgloss.Color("#2F2F2F") // very subtle border
	colorOnline   = lipgloss.Color("#4ADE80") // soft green (not neon)
	colorOffDot   = lipgloss.Color("#F87171") // soft red
)

func (m *model) View() string {
	if m.width == 0 {
		return "Loading…"
	}

	leftW := m.leftW
	rightW := max(10, m.width-leftW-4)

	bodyH := m.height - 4
	bodyH = max(2, bodyH)

	innerH := bodyH - 2
	innerH = max(1, innerH)

	listInnerW := leftW - 4
	if listInnerW < 8 {
		listInnerW = 8
	}

	leftContent := m.renderVehicles(listInnerW, innerH)
	leftBox := lipgloss.NewStyle().
		Width(leftW).
		Height(bodyH).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1).
		Render(leftContent)

	metricsInnerW := max(6, rightW-4)

	rightContent := m.renderMetrics(metricsInnerW)
	rightBox := lipgloss.NewStyle().
		Width(rightW).
		Height(bodyH).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(0, 1).
		Render(rightContent)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	footer := lipgloss.NewStyle().Foreground(colorSubtle).Render(fmt.Sprintf(
		"↑/↓ move · p ping · s shutdown · r apply rate · +/- rate (%dms) · q quit",
		m.rateMs,
	))
	footer = lipgloss.NewStyle().Width(m.width).Render(footer)

	statusLine := ""
	if m.status != "" {
		statusLine = lipgloss.NewStyle().
			Foreground(colorSelected).
			Width(m.width).
			Render(m.status)
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, footer, statusLine)
}

func (m *model) renderVehicles(innerW, innerH int) string {
	if len(m.vehicles) == 0 {
		return lipgloss.NewStyle().Foreground(colorSubtle).Render("Waiting for vehicles…")
	}

	h := m.listHeight()
	if h > innerH {
		h = innerH
	}

	start := m.scroll
	end := start + h
	if end > len(m.vehicles) {
		end = len(m.vehicles)
	}

	lineStyle := lipgloss.NewStyle()
	selStyle := lipgloss.NewStyle().Foreground(colorSelected).Bold(true)
	offStyle := lipgloss.NewStyle().Foreground(colorOffline)

	scrollable := len(m.vehicles) > h && innerW > 4

	rowW := innerW
	if scrollable {
		rowW = innerW - 1
	}
	if rowW < 4 {
		rowW = 4
	}

	rows := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		v := m.vehicles[i]
		on := "offline"
		if v.Online {
			on = "online"
		}
		line := fmt.Sprintf(
			"#%-4d %-9.4f %-9.4f %6.0fkm/h %4d%% %5.0f°C %-7s",
			v.VehicleID,
			v.Lat,
			v.Lon,
			v.SpeedKmh,
			v.BatteryPercent,
			v.EngineTempC,
			on,
		)

		base := lipgloss.NewStyle().Width(rowW - 2) // leave space for "● "
		content := base.Render(util.Truncate(line, rowW-2))

		styled := content
		if i == m.cursor {
			styled = selStyle.Render(content)
		} else if !v.Online {
			styled = offStyle.Render(content)
		} else {
			styled = lineStyle.Render(content)
		}

		dotColor := colorOffDot
		if v.Online {
			dotColor = colorOnline
		}
		dot := lipgloss.NewStyle().Foreground(dotColor).Render("● ")

		row := dot + styled

		rows = append(rows, row)
	}

	listText := strings.Join(rows, "\n")

	if !scrollable {
		return listText
	}

	bar := fleetScrollBarView(len(m.vehicles), m.scroll, h)
	return lipgloss.JoinHorizontal(lipgloss.Top, listText, bar)
}

func (m *model) renderMetrics(innerW int) string {
	_ = innerW

	k := lipgloss.NewStyle().Foreground(colorMetricK)
	v := lipgloss.NewStyle()

	if m.metrics == nil {
		return lipgloss.NewStyle().Foreground(colorSubtle).Render("Waiting for metrics…")
	}
	mm := m.metrics

	online := int(mm.ActiveVehicles)
	offline := int(mm.OfflineVehicles)
	total := online + offline

	barW := innerW - 2
	if barW < 4 {
		barW = 4
	}
	onlineBar := ""
	offlineBar := ""
	if total > 0 && barW > 0 {
		filledN := online * barW / total
		onlineBar = strings.Repeat("█", filledN)
		offlineBar = strings.Repeat("░", barW-filledN)
	}
	barRendered := lipgloss.NewStyle().Foreground(colorOnline).Render(onlineBar) +
		lipgloss.NewStyle().Foreground(colorOffDot).Render(offlineBar)

	ts := time.UnixMilli(mm.TimestampMs).Format("15:04:05")

	var b strings.Builder
	b.WriteString(k.Render("Active   ") + " " + v.Render(fmt.Sprintf("%d", mm.ActiveVehicles)) + "\n")
	b.WriteString(k.Render("Offline  ") + " " + v.Render(fmt.Sprintf("%d", mm.OfflineVehicles)) + "\n")
	b.WriteString(k.Render("Low bat  ") + " " + v.Render(fmt.Sprintf("%d", mm.LowBatteryVehicles)) + "\n")
	b.WriteString(k.Render("Avg spd  ") + " " + v.Render(fmt.Sprintf("%.1f km/h", mm.AvgSpeedKmh)) + "\n")
	b.WriteString(k.Render("Avg eng  ") + " " + v.Render(fmt.Sprintf("%.1f °C", mm.AvgEngineTempC)) + "\n")
	b.WriteString(k.Render("Updated  ") + " " + v.Render(ts) + "\n")
	b.WriteString("\n")
	b.WriteString(barRendered + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(colorSubtle).
		Render(fmt.Sprintf("%d/%d online", online, total)) + "\n")

	if id, ok := m.selectedID(); ok {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(colorSelected).Bold(true).
			Render(fmt.Sprintf("▶ Vehicle #%d", id)))

		for _, veh := range m.vehicles {
			if veh.VehicleID != id {
				continue
			}
			b.WriteString("\n")
			onStr := lipgloss.NewStyle().Foreground(colorOffDot).Render("offline")
			if veh.Online {
				onStr = lipgloss.NewStyle().Foreground(colorOnline).Render("online")
			}
			b.WriteString(k.Render("  Status ") + " " + onStr + "\n")
			b.WriteString(k.Render("  Lat    ") + " " + v.Render(fmt.Sprintf("%.6f", veh.Lat)) + "\n")
			b.WriteString(k.Render("  Lon    ") + " " + v.Render(fmt.Sprintf("%.6f", veh.Lon)) + "\n")
			b.WriteString(k.Render("  Speed  ") + " " + v.Render(fmt.Sprintf("%.1f km/h", veh.SpeedKmh)) + "\n")
			b.WriteString(k.Render("  Bat    ") + " " + v.Render(fmt.Sprintf("%d%%", veh.BatteryPercent)) + "\n")
			b.WriteString(k.Render("  Temp   ") + " " + v.Render(fmt.Sprintf("%.1f °C", veh.EngineTempC)) + "\n")
			break
		}
	}

	return b.String()
}

func fleetScrollBarView(total, offset, height int) string {
	trackStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#5c5f77"))
	thumbStyle := lipgloss.NewStyle().Foreground(colorSelected)

	if height <= 0 || total <= height {
		lines := make([]string, height)
		for i := range lines {
			lines[i] = trackStyle.Render("▏")
		}
		return strings.Join(lines, "\n")
	}

	thumbH := height * height / total
	if thumbH < 1 {
		thumbH = 1
	}
	maxY := total - height
	pos := 0
	if maxY > 0 {
		pos = offset * (height - thumbH) / maxY
	}
	if pos < 0 {
		pos = 0
	}
	if pos+thumbH > height {
		pos = height - thumbH
	}

	lines := make([]string, height)
	for i := range lines {
		if i >= pos && i < pos+thumbH {
			lines[i] = thumbStyle.Render("█")
		} else {
			lines[i] = trackStyle.Render("▏")
		}
	}
	return strings.Join(lines, "\n")
}
