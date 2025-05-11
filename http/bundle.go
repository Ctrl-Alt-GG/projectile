package http

import (
	"github.com/Ctrl-Alt-GG/projectile/db"
	"github.com/Ctrl-Alt-GG/projectile/model"
	"github.com/gin-gonic/gin"
)

type BundleResp struct {
	GameServers  []model.GameServer `json:"game_servers"`
	Announcement string             `json:"announcement"`
}

func getBundle(ctx *gin.Context) {
	ctx.JSON(200, BundleResp{
		GameServers:  db.GetAll(),
		Announcement: db.LoadAnnouncement(),
	})
}
