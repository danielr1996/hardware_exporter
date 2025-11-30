package main

import (
	"os"
	"path/filepath"
	"strconv"
)

type NicInfo struct {
	Name  string
	MAC   string
	MTU   int
	Speed int64 // Mbit/s if available
}

func collectNics() []NicInfo {
	var result []NicInfo

	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return result
	}
	for _, e := range entries {
		ifname := e.Name()
		mac := readFirstLine(filepath.Join("/sys/class/net", ifname, "address"))
		mtuStr := readFirstLine(filepath.Join("/sys/class/net", ifname, "mtu"))
		mtu, _ := strconv.Atoi(mtuStr)

		speedStr := readFirstLine(filepath.Join("/sys/class/net", ifname, "speed"))
		speed, _ := strconv.ParseInt(speedStr, 10, 64)

		result = append(result, NicInfo{
			Name:  ifname,
			MAC:   mac,
			MTU:   mtu,
			Speed: speed * 1_000_000, // convert Mbit → bit/s
		})
	}
	return result
}
