package main

import (
	"flag"
	"github.com/pelletier/go-toml/v2"
	"io"
	"myvgrest/app"
	"myvgrest/models"
	"os"
	"os/user"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		panic("[isRoot] Unable to get current user " + err.Error())
	}
	return currentUser.Username == "root"
}

func main() {
	if !isRoot() {
		panic("needed to be run from root")
	}
	var configPathFlag string
	flag.StringVar(&configPathFlag, "config", "/etc/mywgrest/config.toml", "path to config")
	flag.Parse()

	file, err := os.Open(configPathFlag)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	bts, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var cfg models.Config
	err = toml.Unmarshal(bts, &cfg)
	if err != nil {
		panic(err)
	}

	a := app.NewApplication(cfg)
	a.Run()
}
