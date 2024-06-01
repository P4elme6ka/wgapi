package wg_control

import (
	"github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"myvgrest/models"
	"myvgrest/mylog"
	"net"
	"os"
	"time"
)

func CreatePeer(deviceName string, selfIp net.IP) (*models.PeerWrapper, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	peerConf := wgtypes.PeerConfig{}

	newPeerId := uuid.New()

	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	device, err := client.Device(deviceName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	//err = request.Apply(&peerConf)
	_, allIpsV4, _ := net.ParseCIDR("0.0.0.0/0")
	_, allIpsV6, _ := net.ParseCIDR("::/0")

	keepalive := time.Duration(float64(time.Second) * 30)

	peerConf.AllowedIPs = []net.IPNet{*allIpsV4, *allIpsV6}
	peerConf.PersistentKeepaliveInterval = &keepalive
	peerConf.Endpoint = &net.UDPAddr{
		IP:   selfIp,
		Port: device.ListenPort,
		Zone: "",
	}
	peerConf.PublicKey = privateKey.PublicKey()

	deviceConf := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			peerConf,
		},
	}

	if err := client.ConfigureDevice(deviceName, deviceConf); err != nil {
		return nil, err
	}

	device, err = client.Device(deviceName)
	if err != nil {
		return nil, err
	}

	var peer wgtypes.Peer
	for _, v := range device.Peers {
		if v.PublicKey == peerConf.PublicKey {
			peer = v
			break
		}
	}

	peerWrapper := &models.PeerWrapper{
		Id:         newPeerId,
		PeerConfig: peerConf,
		Peer:       peer,
		Device:     device,
		PrivateKey: privateKey,
	}

	peerWrapper.StringConf, _ = peerWrapper.GetStringConfig()

	return peerWrapper, nil
}

func UpdatePeer(peer *models.PeerWrapper) error {
	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Device(peer.Device.Name)
	if err != nil {
		return err
	}

	deviceConf := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey:         peer.PrivateKey.PublicKey(),
				ReplaceAllowedIPs: true,
				UpdateOnly:        true,
			},
		},
	}

	if err := client.ConfigureDevice(peer.Device.Name, deviceConf); err != nil {
		return err
	}

	return nil
}

func DeletePeer(peer *models.PeerWrapper) error {
	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Device(peer.Device.Name)
	if err != nil {
		return err
	}

	deviceConf := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey: peer.PrivateKey.PublicKey(),
				Remove:    true,
			},
		},
	}

	if err := client.ConfigureDevice(peer.Device.Name, deviceConf); err != nil {
		return err
	}

	return nil
}

func DeleteUnusedPeers(deviceName string, inactivePeriod time.Duration) error {
	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer client.Close()

	device, err := client.Device(deviceName)
	if err != nil {
		return err
	}

	deviceConf := wgtypes.Config{}

	for _, peer := range device.Peers {
		if !peer.LastHandshakeTime.IsZero() && !peer.LastHandshakeTime.Add(inactivePeriod).After(time.Now()) {
			mylog.GetLogger().Infoln("removing peer: ", peer.PublicKey.String())
			deviceConf.Peers = append(deviceConf.Peers, wgtypes.PeerConfig{
				PublicKey: peer.PublicKey,
				Remove:    true,
			})
		}
	}

	if err := client.ConfigureDevice(deviceName, deviceConf); err != nil {
		return err
	}

	return nil
}
