package mylog

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"myvgrest/models"
	"os"
	"path/filepath"
)

var logger *logrus.Logger = nil

func GetLogger() *logrus.Logger {
	if logger == nil {
		panic("logger not inited")
	}
	return logger
}

func SetupLogger(config models.Config, ctx context.Context) {
	if config.Debug {
		logger = logrus.StandardLogger()
		return
	}

	if err := os.MkdirAll(filepath.Dir(config.LogPath), 0770); err != nil {
		log.Println("error on creating folders: ", err)
		return
	}

	// open a file
	f, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = f.Close()
		}
	}()

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(f)
}
