package utils

import (
	"net"
	"strings"
)

// GetClusterDomain retrieves cluster domain of Kubernetes cluster we are installed on.
// It can be retrieved by looking up CNAME of kubernetes.default.svc and extracting its suffix.
func GetClusterDomain() (string, error) {
	apiSvc := "kubernetes.default.svc"
	cname, err := net.LookupCNAME(apiSvc)
	if err != nil {
		return "", err
	}
	clusterDomain := strings.TrimPrefix(cname, apiSvc+".")
	clusterDomain = strings.TrimSuffix(clusterDomain, ".")
	return clusterDomain, nil
}
