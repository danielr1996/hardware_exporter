package collector

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/danielr1996/hardware_exporter/internal/util"
)

type CPUInfo struct {
	Vendor        string
	ModelName     string
	Architecture  string
	PhysicalCores int
	LogicalCores  int
	Sockets       int

	MinHz float64
	MaxHz float64
}

func readCPUFreq(path string) float64 {
	s := util.ReadFirstLine(path)
	v, _ := strconv.ParseFloat(s, 64)
	return v * 1000 // convert kHz → Hz
}

func CollectCPUInfo() CPUInfo {
	logical := runtime.NumCPU()

	sockets := map[string]bool{}
	cores := map[string]bool{}
	vendor := ""
	model := ""

	f, err := os.Open("/proc/cpuinfo")
	if err == nil {
		defer f.Close()

		sc := bufio.NewScanner(f)
		var physID, coreID string
		for sc.Scan() {
			line := sc.Text()

			if strings.HasPrefix(line, "physical id") {
				physID = strings.TrimSpace(strings.Split(line, ":")[1])
				sockets[physID] = true
			}
			if strings.HasPrefix(line, "core id") {
				coreID = strings.TrimSpace(strings.Split(line, ":")[1])
				cores[physID+"-"+coreID] = true
			}
			if strings.HasPrefix(line, "vendor_id") && vendor == "" {
				vendor = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			if strings.HasPrefix(line, "model name") && model == "" {
				model = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	}

	physical := len(cores)
	if physical == 0 {
		physical = logical
	}

	minHz := readCPUFreq("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_min_freq")
	maxHz := readCPUFreq("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")

	return CPUInfo{
		Vendor:        vendor,
		ModelName:     model,
		Architecture:  runtime.GOARCH,
		PhysicalCores: physical,
		LogicalCores:  logical,
		Sockets:       len(sockets),

		MinHz: minHz,
		MaxHz: maxHz,
	}
}
