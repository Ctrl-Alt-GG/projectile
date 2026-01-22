package http

import (
	"net/http"

	"github.com/Ctrl-Alt-GG/projectile/cmd/server/db"
	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func updateGameServer(ctx *gin.Context) {
	id := ctx.Param("id")
	l := GetLoggerFromContext(ctx).With(zap.String("id", id))
	if id == "" {
		l.Warn("ID is invalid")
		ctx.Status(http.StatusBadRequest)
		return
	}

	var update model.GameServerData
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

	db.StoreGameServerData(id, update) // this updates the last update value
	l.Debug("Update stored!")
	ctx.Status(http.StatusAccepted)
}

func deleteGameServer(ctx *gin.Context) {
	id := ctx.Param("id")
	l := GetLoggerFromContext(ctx).With(zap.String("id", id))
	if id == "" {
		l.Warn("ID is invalid")
		ctx.Status(http.StatusBadRequest)
		return
	}

	cnt := db.DeleteServer(id)
	l.Debug("Record deletion completed", zap.Int("deletedRecords", cnt))

	if cnt == 0 {
		ctx.Status(http.StatusNotFound)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

func getGameServer(ctx *gin.Context) {
	id := ctx.Param("id")
	l := GetLoggerFromContext(ctx).With(zap.String("id", id))
	if id == "" {
		l.Warn("ID is invalid")
		ctx.Status(http.StatusBadRequest)
		return
	}

	srv, ok := db.GetServer(id)
	if ok {
		ctx.JSON(http.StatusOK, srv)
	} else {
		ctx.Status(http.StatusNotFound)
	}
}

func getAllGameServersAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, db.GetMap())
}

func getAllGameServersPublic(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, db.GetList())
}
