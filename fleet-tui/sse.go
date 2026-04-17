package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type vehicleDTO struct {
	VehicleID       uint32  `json:"vehicle_id"`
	Lat             float64 `json:"lat"`
	Lon             float64 `json:"lon"`
	SpeedKmh        float32 `json:"speed_kmh"`
	BatteryPercent  uint32  `json:"battery_percent"`
	EngineTempC     float32 `json:"engine_temp_c"`
	LastTelemetryMs int64   `json:"last_telemetry_ms"`
	Online          bool    `json:"online"`
}

type fleetVehiclesDTO struct {
	Vehicles    []vehicleDTO `json:"vehicles"`
	TimestampMs int64        `json:"timestamp_ms"`
}

type fleetMetricsDTO struct {
	ActiveVehicles     uint32  `json:"active_vehicles"`
	OfflineVehicles    uint32  `json:"offline_vehicles"`
	LowBatteryVehicles uint32  `json:"low_battery_vehicles"`
	AvgSpeedKmh        float32 `json:"avg_speed_kmh"`
	AvgEngineTempC     float32 `json:"avg_engine_temp_c"`
	TimestampMs        int64   `json:"timestamp_ms"`
}

type vehiclesDataMsg struct {
	data fleetVehiclesDTO
}

type metricsDataMsg struct {
	data fleetMetricsDTO
}

type statusMsg struct {
	text string
	err  error
}

func streamVehiclesSSE(p *tea.Program, base string) {
	client := &http.Client{Timeout: 0}
	for {
		req, err := http.NewRequest(http.MethodGet, strings.TrimRight(base, "/")+"/vehicles/stream", nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		readVehicleSSE(resp.Body, p)
		resp.Body.Close()
		time.Sleep(time.Second)
	}
}

// sseScannerMaxLine is the max bytes per SSE line. Default bufio.Scanner limit is 64KiB;
// fleet JSON can exceed that with many vehicles, which made Scan() fail with no visible error.
const sseScannerMaxLine = 8 << 20 // 8 MiB

func readVehicleSSE(body io.Reader, p *tea.Program) {
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64*1024), sseScannerMaxLine)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		var out fleetVehiclesDTO
		if err := json.Unmarshal([]byte(payload), &out); err != nil {
			continue
		}
		p.Send(vehiclesDataMsg{data: out})
	}
	_ = sc.Err()
}

func streamMetricsSSE(p *tea.Program, base string) {
	client := &http.Client{Timeout: 0}
	for {
		req, err := http.NewRequest(http.MethodGet, strings.TrimRight(base, "/")+"/metrics/stream", nil)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		readMetricsSSE(resp.Body, p)
		resp.Body.Close()
		time.Sleep(time.Second)
	}
}

func readMetricsSSE(body io.Reader, p *tea.Program) {
	sc := bufio.NewScanner(body)
	sc.Buffer(make([]byte, 0, 64*1024), sseScannerMaxLine)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		var out fleetMetricsDTO
		if err := json.Unmarshal([]byte(payload), &out); err != nil {
			continue
		}
		p.Send(metricsDataMsg{data: out})
	}
	_ = sc.Err()
}

func postCommand(base string, vehicleID uint32, typ string, value uint32) error {
	body := map[string]any{
		"type":  typ,
		"value": value,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/vehicles/%d/command", strings.TrimRight(base, "/"), vehicleID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slurp, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(slurp)))
	}
	return nil
}
