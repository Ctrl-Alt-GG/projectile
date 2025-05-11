package http

import (
	"github.com/Ctrl-Alt-GG/projectile/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAnnouncement(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, db.LoadAnnouncement())
}

func setAnnouncement(ctx *gin.Context) {
	var newAnnouncement string
	err := ctx.BindJSON(&newAnnouncement)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	err = db.StoreAnnouncement(newAnnouncement)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func clearAnnouncement(ctx *gin.Context) {
	err := db.StoreAnnouncement("")
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
