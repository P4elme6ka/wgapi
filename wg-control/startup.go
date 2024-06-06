package wg_control

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"myvgrest/mylog"
	"myvgrest/utils"
)

func CheckWgInstallation() {
	wgclient, err := wgctrl.New()
	if err != nil {
		mylog.GetLogger().Fatalf("failed to open wgctrl: %v", err)
	}

	myip, err := utils.Getip2()
	mylog.GetLogger().Infof("current external ip = %s\n", myip.String())

	defer func() {
		err := wgclient.Close()
		mylog.GetLogger().Panicln(err)
	}()
}
