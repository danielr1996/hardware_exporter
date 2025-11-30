package collector

import (
	"github.com/danielr1996/hardware_exporter/internal/util"
)

type BaseboardInfo struct {
	Vendor  string
	Name    string
	Version string
	Serial  string
}

func CollectBaseboardInfo() BaseboardInfo {
	return BaseboardInfo{
		Vendor:  util.ReadFirstLine("/sys/class/dmi/id/board_vendor"),
		Name:    util.ReadFirstLine("/sys/class/dmi/id/board_name"),
		Version: util.ReadFirstLine("/sys/class/dmi/id/board_version"),
		Serial:  util.ReadFirstLine("/sys/class/dmi/id/board_serial"),
	}
}
