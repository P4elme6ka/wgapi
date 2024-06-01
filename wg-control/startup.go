package wg_control

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"myvgrest/mylog"
)

func CheckWgInstallation() {
	wgclient, err := wgctrl.New()
	if err != nil {
		mylog.GetLogger().Fatalf("failed to open wgctrl: %v", err)
	}

	defer func() {
		err := wgclient.Close()
		mylog.GetLogger().Panicln(err)
	}()
}
