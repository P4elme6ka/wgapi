package utils

import (
	externalip "github.com/glendc/go-external-ip"
	"net"
)

func GetExternalIP() (net.IP, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return nil, err
	}

	return ip, nil
}
