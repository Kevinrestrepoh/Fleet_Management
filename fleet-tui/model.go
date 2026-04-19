package main

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	baseURL string
	prog    *tea.Program

	vehicles    []vehicleDTO
	metrics     *fleetMetricsDTO
	status      string
	cursor      int
	scroll      int
	rateMs      uint32
	width       int
	height      int
	innerH      int
	leftW       int
	headerLines int
}

func initialModel(base string) *model {
	return &model{
		baseURL: base,
		rateMs:  500,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}
