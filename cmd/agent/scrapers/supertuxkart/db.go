package supertuxkart

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ctrl-Alt-GG/projectile/pkg/model"
	"go.uber.org/zap"
)

const queryStmt = `SELECT username, player_num FROM v1_server_config_current_players;`

func QueryPlayers(ctx context.Context, logger *zap.Logger, dbPath string) ([]model.Player, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", dbPath))
	if err != nil {
		logger.Error("Failed to open STK sqlite db", zap.Error(err))
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Warn("Failed to close the sqlite db", zap.Error(err))
		}
	}(db)

	rows, err := db.QueryContext(ctx, queryStmt)
	if err != nil {
		logger.Error("Failed to query the STK sqlite db", zap.Error(err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Warn("Failed to close rows", zap.Error(err))
		}
	}(rows)

	players := make([]model.Player, 0)

	for rows.Next() {
		var username string
		var playerNum int
		err = rows.Scan(&username, &playerNum)
		if err != nil {
			logger.Error("Failure while scanning rows", zap.Error(err))
			return nil, err
		}

		for i := range playerNum {
			localUsername := username
			if i > 0 {
				localUsername = fmt.Sprintf("%s (%d)", localUsername, i+1)
			}
			players = append(players, model.Player{
				Name: localUsername,
			})
		}
	}

	return players, nil
}
