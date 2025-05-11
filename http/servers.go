package http

import (
	"github.com/Ctrl-Alt-GG/projectile/db"
	"github.com/Ctrl-Alt-GG/projectile/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func updateGameServer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)

	var update model.GameServer
	err := ctx.BindJSON(&update)
	if err != nil {
		l.Warn("Unable to parse body", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = update.Validate()
	if err != nil {
		l.Warn("Body has invalid values", zap.Error(err))
		ctx.Status(http.StatusUnprocessableEntity)
		return
	}

	update.LastUpdate = time.Now()

	err = db.StoreUpdate(update)
	if err != nil {
		l.Error("Failed to store update", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
	}

	l.Debug("Update stored!")

	ctx.Status(http.StatusAccepted)
}

func deleteGameServer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)

	addr := ctx.Param("addr")
	if addr == "" {
		l.Warn("Can not delete empty address")
		ctx.Status(http.StatusBadRequest)
	}

	cnt, err := db.DeleteServer(addr)
	if err != nil {
		l.Error("Failed to delete record", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
	}

	l.Debug("Record deletion completed", zap.Int("deletedRecords", cnt))

	if cnt == 0 {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

func getAllGameServers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, db.GetAll())
}
