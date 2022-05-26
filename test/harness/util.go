package harness

import (
	"net"
	"os/exec"
)

// LocalBinAvailable returns true if the binary is available on PATH.
func LocalBinAvailable(path string) (string, error) {
	return exec.LookPath(path)
}

// AllocateLocalAddress listens on the given address and returns it.
func AllocateLocalAddress(address string) (string, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}
