package http

import (
	"net/http"
	"time"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/db"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		return
	}

	l.Debug("Update stored!")

	ctx.Status(http.StatusAccepted)
}

func deleteGameServer(ctx *gin.Context) {
	l := GetLoggerFromContext(ctx)

	id := model.Identifier{
		Address: ctx.Param("addr"),
		Game:    ctx.Param("game"),
	}

	if !id.IsValid() {
		l.Warn("ID is invalid", zap.String("id.Address", id.Address), zap.String("id.Game", id.Game))
		ctx.Status(http.StatusBadRequest)
		return
	}

	cnt, err := db.DeleteServer(id)
	if err != nil {
		l.Error("Failed to delete record", zap.Error(err))
		ctx.Status(http.StatusInternalServerError)
		return
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
