package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerKey = "lgr"
)

func goodLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		// some evil middlewares may modify this value, so we store it
		path := ctx.Request.URL.Path

		subLogger := logger.With(
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
		)

		ctx.Set(loggerKey, subLogger)

		ctx.Next() // <- execute next thing in the chain
		end := time.Now()

		latency := end.Sub(start)

		completedRequestFields := []zapcore.Field{
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", latency),
		}

		if len(ctx.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range ctx.Errors.Errors() {
				subLogger.Error(e, completedRequestFields...)
			}
		}

		subLogger.Info(fmt.Sprintf("%s %s served: %d", ctx.Request.Method, path, ctx.Writer.Status()), completedRequestFields...) // <- always print this
	}
}

func GetLoggerFromContext(ctx *gin.Context) *zap.Logger { // This one panics
	var logger *zap.Logger
	l, ok := ctx.Get(loggerKey)
	if !ok {
		panic("logger was not in context")
	}
	logger = l.(*zap.Logger)
	return logger
}

func validateAuthHeader(ctx *gin.Context, key string) bool {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return false
	}

	if parts[0] != "Key" {
		return false
	}

	if parts[1] != key {
		return false
	}

	return true
}

func goodKeyAuthMiddleware(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if validateAuthHeader(ctx, key) {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func goodCORSMiddleware(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
