package operations

import (
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/sirupsen/logrus"
)

// HardwareService provides system hardware information
type HardwareService struct {
	logger *logrus.Logger
}

// NewHardwareService creates a new hardware service
func NewHardwareService(logger *logrus.Logger) *HardwareService {
	return &HardwareService{
		logger: logger,
	}
}

// DiskInfo represents physical disk usage
type DiskInfo struct {
	MountPoint  string  `json:"mount_point"`
	Filesystem  string  `json:"filesystem"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkInterfaceInfo represents network interface details
type NetworkInterfaceInfo struct {
	Name        string   `json:"name"`
	MAC         string   `json:"mac"`
	IPs         []string `json:"ips"`
	Flags       []string `json:"flags"`
	BytesSent   uint64   `json:"bytes_sent"`
	BytesRecv   uint64   `json:"bytes_recv"`
	PacketsSent uint64   `json:"packets_sent"`
	PacketsRecv uint64   `json:"packets_recv"`
}

// GetStorageInfo returns information about physical storage
func (s *HardwareService) GetStorageInfo() ([]DiskInfo, error) {
	var disks []DiskInfo

	// Get partitions (only physical ones)
	partitions, err := disk.Partitions(false) // false = all, not just physical. Wait, doc says 'false' for all?
	// godoc: Partitions(all bool) ([]PartitionStat, error) . If all is false, it returns physical devices only (e.g. hard disks, cd-rom drives, USB keys).
	if err != nil {
		s.logger.WithError(err).Error("Failed to get partitions")
		return nil, err
	}

	seen := make(map[string]bool)

	for _, p := range partitions {
		// Filter out snap, docker, loops, and special filesystems if unwanted
		if strings.HasPrefix(p.Mountpoint, "/var/lib/docker") ||
			strings.HasPrefix(p.Mountpoint, "/run") ||
			strings.HasPrefix(p.Mountpoint, "/sys") ||
			strings.HasPrefix(p.Mountpoint, "/proc") ||
			strings.HasPrefix(p.Mountpoint, "/dev") {
			continue
		}

		if seen[p.Mountpoint] {
			continue
		}
		seen[p.Mountpoint] = true

		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			s.logger.WithField("mount", p.Mountpoint).Warn("Failed to get disk usage, skipping")
			continue
		}

		disks = append(disks, DiskInfo{
			MountPoint:  p.Mountpoint,
			Filesystem:  p.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	// Always ensure root '/' is present if missed
	if !seen["/"] {
		usage, err := disk.Usage("/")
		if err == nil {
			disks = append(disks, DiskInfo{
				MountPoint:  "/",
				Filesystem:  "rootfs",
				Total:       usage.Total,
				Free:        usage.Free,
				Used:        usage.Used,
				UsedPercent: usage.UsedPercent,
			})
		}
	}

	return disks, nil
}

// GetNetworkInfo returns information about network interfaces
func (s *HardwareService) GetNetworkInfo() ([]NetworkInterfaceInfo, error) {
	var results []NetworkInterfaceInfo

	interfaces, err := net.Interfaces()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get network interfaces")
		return nil, err
	}

	ioCounters, err := net.IOCounters(true)
	ioMap := make(map[string]net.IOCountersStat)
	if err == nil {
		for _, stat := range ioCounters {
			ioMap[stat.Name] = stat
		}
	}

	for _, iface := range interfaces {
		// Filter loopback if desired, but user might want to see it. Keeping it for now.
		var info NetworkInterfaceInfo
		info.Name = iface.Name
		info.MAC = iface.HardwareAddr
		info.Flags = iface.Flags

		for _, addr := range iface.Addrs {
			info.IPs = append(info.IPs, addr.Addr)
		}

		if stat, ok := ioMap[iface.Name]; ok {
			info.BytesSent = stat.BytesSent
			info.BytesRecv = stat.BytesRecv
			info.PacketsSent = stat.PacketsSent
			info.PacketsRecv = stat.PacketsRecv
		}

		results = append(results, info)
	}

	return results, nil
}
