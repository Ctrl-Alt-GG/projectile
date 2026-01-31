package satisfactory

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type QueryResponse struct {
	Data struct {
		ServerGameState struct {
			ActiveSessionName   string  `json:"activeSessionName"`
			NumConnectedPlayers int     `json:"numConnectedPlayers"`
			PlayerLimit         int     `json:"playerLimit"`
			TechTier            int     `json:"techTier"`
			ActiveSchematic     string  `json:"activeSchematic"`
			GamePhase           string  `json:"gamePhase"`
			IsGameRunning       bool    `json:"isGameRunning"`
			TotalGameDuration   int     `json:"totalGameDuration"`
			IsGamePaused        bool    `json:"isGamePaused"`
			AverageTickRate     float64 `json:"averageTickRate"`
			AutoLoadSessionName string  `json:"autoLoadSessionName"`
		} `json:"serverGameState"`
	} `json:"data"`
}

var insecureClient *http.Client

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // TODO: we could verify the self-signed cert, but I'm lazy
	}
	insecureClient = &http.Client{Transport: tr}
}

func doQuery(ctx context.Context, logger *zap.Logger, url, token string) (QueryResponse, error) {
	jsonStr := []byte(`{"function":"QueryServerState"}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Error("Failed to craft new request", zap.Error(err))
		return QueryResponse{}, err
	}
	req.Header.Set("Content-type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := insecureClient.Do(req)
	if err != nil {
		logger.Error("Failed to do HTTP request", zap.Error(err))
		return QueryResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warn("Failed to close body", zap.Error(err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected http status: %d %s", resp.StatusCode, resp.Status)
		logger.Error("Unexpected http status", zap.Error(err))
		return QueryResponse{}, err
	}

	var queryResp QueryResponse
	err = json.NewDecoder(resp.Body).Decode(&queryResp)
	if err != nil {
		logger.Error("Failed to parse response", zap.Error(err))
		return QueryResponse{}, err
	}

	return queryResp, nil
}
