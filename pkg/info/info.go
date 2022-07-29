package info

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	infov1 "github.com/FluxNinja/aperture/api/gen/proto/go/aperture/common/info/v1"
)

// Default build-time variables. These values are overridden via ldflags.
var (
	// Version is the version of the aperture process.
	Version = "0.0.1"
	// Service is the name of the aperture process.
	Service = "unknown"
	// BuildHost is the hostname of the machine that built the application.
	BuildHost = "unknown"
	// BuildOS is the operating system that built the application.
	BuildOS = "unknown"
	// BuildTime is the time when the application was built.
	BuildTime = "unknown"
	// GitBranch is the git branch that the application was built from.
	GitBranch = "unknown"
	// GitCommitHash is the git commit hash that the application was built from.
	GitCommitHash = "unknown"
	// Prefix is the prefix for the aperture service application.
	Prefix = "aperture"
	// Hostname is the hostname of the machine that is running the process.
	Hostname = "unknown"
)

var (
	mutex       sync.Mutex
	versionInfo infov1.VersionInfo
	processInfo infov1.ProcessInfo
)

func init() {
	if Service == "unknown" {
		service := filepath.Base(os.Args[0])
		// remove dots from service name if it has them
		Service = strings.ReplaceAll(service, ".", "-")
	}
	versionInfo.Version = Version
	versionInfo.Service = Service
	versionInfo.BuildHost = BuildHost
	versionInfo.BuildOs = BuildOS
	versionInfo.BuildTime = BuildTime
	versionInfo.GitBranch = GitBranch
	versionInfo.GitCommitHash = GitCommitHash

	// ProcessInfo
	processInfo.StartTime = timestamppb.Now()

	// Hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = getLocalIP()
	}
	if hostname != "" {
		Hostname = hostname
	}
}

// GetVersionInfo returns the version info for the application.
func GetVersionInfo() *infov1.VersionInfo {
	mutex.Lock()
	defer mutex.Unlock()
	return proto.Clone(&versionInfo).(*infov1.VersionInfo)
}

// GetProcessInfo returns the process info for the application.
func GetProcessInfo() *infov1.ProcessInfo {
	mutex.Lock()
	defer mutex.Unlock()
	// reset uptime
	processInfo.Uptime = durationpb.New(time.Since(processInfo.StartTime.AsTime()))
	return (proto.Clone(&processInfo)).(*infov1.ProcessInfo)
}

// GetHostInfo returns the host info for the application.
func GetHostInfo() *infov1.HostInfo {
	mutex.Lock()
	defer mutex.Unlock()
	return &infov1.HostInfo{
		Hostname: Hostname,
	}
}

// getLocalIP gets local ip address.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
