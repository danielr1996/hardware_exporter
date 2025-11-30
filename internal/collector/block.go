package collector

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danielr1996/hardware_exporter/internal/util"
)

type BlockDevice struct {
	Name       string
	SizeBytes  uint64
	Rotational bool
}

func CollectBlockDevices() []BlockDevice {
	var result []BlockDevice

	entries, err := os.ReadDir("/sys/block")
	if err != nil {
		log.Printf("read /sys/block failed: %v", err)
		return result
	}

	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "loop") ||
			strings.HasPrefix(name, "dm-") ||
			strings.HasPrefix(name, "zram") ||
			strings.HasPrefix(name, "ram") ||
			strings.HasPrefix(name, "nbd") {
			continue
		}

		// Real path (follow symlink)
		devicePath := filepath.Join("/sys/block", name)

		// size is inside the symlink target dir
		sizeStr := util.ReadFirstLine(filepath.Join(devicePath, "size"))
		rotStr := util.ReadFirstLine(filepath.Join(devicePath, "queue/rotational"))

		secs, _ := strconv.ParseUint(sizeStr, 10, 64)
		size := secs * 512
		rot := rotStr == "1"

		result = append(result, BlockDevice{
			Name:       name,
			SizeBytes:  size,
			Rotational: rot,
		})
	}

	return result
}
