package utils

import (
	"encoding/json"
	externalip "github.com/glendc/go-external-ip"
	"io/ioutil"
	"net"
	"net/http"
)

type IP struct {
	Query string
}

func Getip2() (net.IP, error) {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return net.ParseIP(ip.Query), nil
}

func GetExternalIP() (net.IP, error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		return nil, err
	}

	return ip, nil
}
