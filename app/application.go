package app

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"myvgrest/handlers"
	"myvgrest/models"
	"myvgrest/mylog"
	"myvgrest/store"
	"myvgrest/utils"
	wg_control "myvgrest/wg-control"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	config models.Config
	router *gin.Engine
	ctx    context.Context
}

func NewApplication(config models.Config) *Application {
	context, _ := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	mylog.SetupLogger(config, context)

	router := gin.New()
	router.Use(mylog.Logger(mylog.GetLogger()), gin.Recovery())
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	storage, err := store.OpenStorage(config)
	if err != nil {
		panic(err)
	}

	publicIp, err := utils.GetExternalIP()
	if err != nil {
		panic(err)
	}
	if config.InactivePeersDelete {
		go func() {
			ticker := time.NewTicker(time.Minute)
			for {
				select {
				case <-ticker.C:
					err := wg_control.DeleteUnusedPeers(config.DefaultWgDevice, time.Minute*5)
					if err != nil {
						mylog.GetLogger().Error(err)
					}
				case <-context.Done():
					mylog.GetLogger().Println("exiting cleanup loop")
					return
				}

			}
		}()
	}

	router.GET("new", handlers.AddNewPeer(config.DefaultWgDevice, publicIp, storage))
	//router.GET("get", handlers.GetPeer(storage))
	router.GET("delete", handlers.RemovePeer(storage))
	//router.GET("list", handlers.GetPeerList(storage))

	return &Application{
		config: config,
		router: router,
		ctx:    context,
	}
}

func (a *Application) Run() {
	go func() {
		err := a.router.Run(a.config.ListenAddr)
		if err != nil {
			mylog.GetLogger().Fatalln(err)
		}
	}()

	<-a.ctx.Done()
	return
}
