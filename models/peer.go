package models

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"io"
	"strings"
)

type PeerWrapper struct {
	Id         uuid.UUID
	PeerConfig wgtypes.PeerConfig
	Peer       wgtypes.Peer
	Device     *wgtypes.Device
	PrivateKey wgtypes.Key
	StringConf string
}

func (pw *PeerWrapper) GetStringConfig() (string, error) {
	b := &bytes.Buffer{}
	fmt.Fprintln(b, "[Interface]")
	fmt.Fprintln(b, "PrivateKey =", pw.PrivateKey)

	//addresses := make([]string, len(pw.PeerConfig.AllowedIPs))
	//for i, v := range pw.PeerConfig.AllowedIPs {
	//	addresses[i] = v.String()
	//}
	//
	//fmt.Fprintf(b, "Address = %s\n", strings.Join(addresses, ","))
	fmt.Fprintf(b, "Address = 10.10.3.2\n")
	//if pw.Device.DNSServers != nil && len(*options.DNSServers) > 0 {
	//	fmt.Fprintf(b, "DNS = %s\n", strings.Join(*options.DNSServers, ","))
	//}

	fmt.Fprintln(b, "")
	fmt.Fprintln(b, "[Peer]")

	fmt.Fprintf(b, "PublicKey = %s\n", pw.Device.PublicKey.String())
	emptyKey := wgtypes.Key{}
	if pw.Peer.PresharedKey != emptyKey {
		fmt.Fprintf(b, "PresharedKey = %s\n", pw.Peer.PresharedKey.String())
	}
	if pw.Peer.Endpoint != nil {
		fmt.Fprintf(b, "Endpoint = %s:%d\n", pw.Peer.Endpoint.IP.String(), pw.Peer.Endpoint.Port)
	}
	allowedIps := []string{}
	for _, ip := range pw.PeerConfig.AllowedIPs {
		allowedIps = append(allowedIps, ip.String())
	}
	if pw.PeerConfig.AllowedIPs != nil && len(pw.PeerConfig.AllowedIPs) > 0 {
		fmt.Fprintf(b, "AllowedIPs = %s\n", strings.Join(allowedIps, ","))
	}

	resultString, _ := io.ReadAll(b)

	return string(resultString), nil
}
