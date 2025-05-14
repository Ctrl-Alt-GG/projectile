package http

import (
	"github.com/Ctrl-Alt-GG/projectile/db"
	"github.com/Ctrl-Alt-GG/projectile/model"
	"github.com/gin-gonic/gin"
)

type BundleResp struct {
	Announcement AnnouncementWrapper `json:"announcement"`
	GameServers  []model.GameServer  `json:"gameServers"`
}

func getBundle(ctx *gin.Context) {
	ctx.JSON(200, BundleResp{
		GameServers:  db.GetAll(),
		Announcement: AnnouncementWrapper{db.LoadAnnouncement()},
	})
}
