package http

import (
	"net/http"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/db"
	"github.com/gin-gonic/gin"
)

// It will be easier to deal with like this
type AnnouncementWrapper struct {
	Text string `json:"text"`
}

func getAnnouncement(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, AnnouncementWrapper{db.LoadAnnouncement()})
}

func setAnnouncement(ctx *gin.Context) {
	var newAnnouncement AnnouncementWrapper
	err := ctx.BindJSON(&newAnnouncement)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if newAnnouncement.Text == "" {
		// don't allow empty
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	err = db.StoreAnnouncement(newAnnouncement.Text)
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
