package main

import (
	"os"
	"path/filepath"
	"strconv"
)

type BlockDevice struct {
	Name       string
	SizeBytes  uint64
	Rotational bool
}

func collectBlockDevices() []BlockDevice {
	var result []BlockDevice

	filepath.WalkDir("/sys/block", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if path == "/sys/block" {
			return nil
		}
		if d.IsDir() && filepath.Dir(path) == "/sys/block" {
			name := d.Name()

			sizeStr := readFirstLine(filepath.Join(path, "size"))
			rotStr := readFirstLine(filepath.Join(path, "queue/rotational"))

			secs, _ := strconv.ParseUint(sizeStr, 10, 64)
			size := secs * 512

			rot := rotStr == "1"

			result = append(result, BlockDevice{
				Name:       name,
				SizeBytes:  size,
				Rotational: rot,
			})
		}
		return nil
	})

	return result
}
