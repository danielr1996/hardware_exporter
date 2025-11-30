package main

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type MemoryModule struct {
	Locator       string
	Bank          string
	Channel       string
	Type          string
	Form          string
	Speed         string
	SpeedHz       float64
	ConfigSpeed   string
	ConfigSpeedHz float64
	Manufacturer  string
	Serial        string
	Part          string
	Rank          string
	RankN         float64
	SizeBytes     uint64
	CLLatency     string
	CLCycles      float64
}

func parseSize(s string) uint64 {
	s = strings.ToUpper(strings.TrimSpace(s))
	if strings.HasSuffix(s, "GB") {
		n := strings.TrimSuffix(s, "GB")
		v, _ := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return uint64(v * 1024 * 1024 * 1024)
	}
	if strings.HasSuffix(s, "MB") {
		n := strings.TrimSuffix(s, "MB")
		v, _ := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return uint64(v * 1024 * 1024)
	}
	return 0
}

func parseSpeedHz(s string) float64 {
	// "1600 MT/s" → 1600e6
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return 0
	}
	n, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	return n * 1_000_000
}

func extractChannel(locator string) string {
	// ChannelA-DIMM0 → "A"
	if strings.Contains(locator, "Channel") {
		idx := strings.Index(locator, "Channel")
		if idx >= 0 && len(locator) > idx+7 {
			return string(locator[idx+7])
		}
	}
	return ""
}

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func collectMemory() []MemoryModule {
	out, err := exec.Command("dmidecode", "-t", "memory").Output()
	if err != nil {
		return nil
	}

	var modules []MemoryModule
	var cur MemoryModule
	sc := bufio.NewScanner(strings.NewReader(string(out)))

	reKV := regexp.MustCompile(`^\s*([^:]+):\s*(.*)$`)

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "Memory Device") {
			if cur.Locator != "" {
				modules = append(modules, cur)
			}
			cur = MemoryModule{}
			continue
		}

		m := reKV.FindStringSubmatch(line)
		if len(m) != 3 {
			continue
		}

		key := strings.TrimSpace(m[1])
		val := strings.TrimSpace(m[2])

		switch key {
		case "Locator":
			cur.Locator = val
			cur.Channel = extractChannel(val)
		case "Bank Locator":
			cur.Bank = val
		case "Type":
			cur.Type = val
		case "Form Factor":
			cur.Form = val
		case "Speed":
			cur.Speed = val
			cur.SpeedHz = parseSpeedHz(val)
		case "Configured Memory Speed":
			cur.ConfigSpeed = val
			cur.ConfigSpeedHz = parseSpeedHz(val)
		case "Manufacturer":
			cur.Manufacturer = val
		case "Serial Number":
			cur.Serial = val
		case "Part Number":
			cur.Part = val
		case "Rank":
			cur.Rank = val
			cur.RankN = toFloat(val)
		case "Size":
			cur.SizeBytes = parseSize(val)
		case "CL Latency":
			cur.CLLatency = val
			cur.CLCycles = toFloat(val)
		}
	}

	if cur.Locator != "" {
		modules = append(modules, cur)
	}

	return modules
}
