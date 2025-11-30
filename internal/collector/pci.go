package collector

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/danielr1996/hardware_exporter/internal/util"
)

type PCIDevice struct {
	Address    string
	VendorID   string
	DeviceID   string
	Class      string
	Driver     string
	DeviceName string
}

func lookupPCIName(addr string) string {
	out, err := exec.Command("lspci", "-s", addr).Output()
	if err != nil {
		return ""
	}
	parts := strings.SplitN(string(out), ":", 3)
	if len(parts) < 3 {
		return strings.TrimSpace(string(out))
	}
	return strings.TrimSpace(parts[2])
}

func CollectPCIDevices() []PCIDevice {
	root := "/sys/bus/pci/devices"
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	var out []PCIDevice

	for _, e := range entries {
		addr := e.Name()
		base := filepath.Join(root, addr)

		vendor := strings.TrimPrefix(util.ReadFirstLine(filepath.Join(base, "vendor")), "0x")
		device := strings.TrimPrefix(util.ReadFirstLine(filepath.Join(base, "device")), "0x")
		class := strings.TrimPrefix(util.ReadFirstLine(filepath.Join(base, "class")), "0x")

		driver := ""
		if t, err := os.Readlink(filepath.Join(base, "driver")); err == nil {
			driver = filepath.Base(t)
		}

		name := lookupPCIName(addr)

		out = append(out, PCIDevice{
			Address:    addr,
			VendorID:   vendor,
			DeviceID:   device,
			Class:      class,
			Driver:     driver,
			DeviceName: name,
		})
	}

	return out
}
