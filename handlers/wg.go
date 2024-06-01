package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"myvgrest/mylog"
	"myvgrest/store"
	wg_control "myvgrest/wg-control"
	"net"
	"net/http"
)

func AddNewPeer(deviceName string, publicIp net.IP, storage *store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		peerWrap, err := wg_control.CreatePeer(deviceName, publicIp)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusBadRequest, "failed parse peer id")
			return
		}

		err = storage.SetPeer(peerWrap)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusInternalServerError, "failed write to store")
			return
		}

		GenerateResponse(c, http.StatusOK, peerWrap.StringConf)
	}
}

func RemovePeer(storage *store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		peerIds, ok := c.GetQuery("peerId")
		if !ok || peerIds == "" {
			GenerateMessage(c, http.StatusBadRequest, "failed: peer id is empty")
			return
		}
		peerId, err := uuid.Parse(peerIds)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusBadRequest, "failed: parse peer id")
			return
		}
		peerWrap, err := storage.GetPeer(peerId)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusBadRequest, "failed: peer not found")
			return
		}

		err = wg_control.DeletePeer(peerWrap)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusBadRequest, "failed: delete peer")
			return
		}

		GenerateMessage(c, http.StatusOK, "successfully deleted")
	}
}

func GetPeer(storage *store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		peerIds, ok := c.GetQuery("peerId")
		if !ok || peerIds == "" {
			GenerateMessage(c, http.StatusBadRequest, "failed: peer id is empty")
			return
		}
		peerId, err := uuid.Parse(peerIds)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusBadRequest, "failed: parse peer id")
			return
		}
		peerWrap, err := storage.GetPeer(peerId)
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusNotFound, "failed: peer not found")
			return
		}

		GenerateResponse(c, http.StatusOK, peerWrap.StringConf)
	}
}

func GetPeerList(storage *store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		peers, err := storage.ListPeer()
		if err != nil {
			mylog.GetLogger().Error(err)
			GenerateMessage(c, http.StatusInternalServerError, "failed: peer not found")
			return
		}
		GenerateResponse(c, http.StatusOK, peers)
	}
}

func GenerateMessage(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"message": msg,
		"data":    nil,
		"code":    code,
	})
}

func GenerateResponse(ctx *gin.Context, code int, payload interface{}) {
	ctx.JSON(code, gin.H{
		"message": "",
		"data":    payload,
		"code":    code,
	})
}
