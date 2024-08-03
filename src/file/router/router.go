package router

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/file/global"
	"github.com/whoisnian/k8s-example/src/file/pkg/apis"
	"github.com/whoisnian/k8s-example/src/file/router/file"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Setup() *gin.Engine {
	if global.CFG.Debug {
		gin.SetMode(gin.DebugMode)
		gin.DebugPrintFunc = func(format string, values ...interface{}) {
			global.LOG.Debug(fmt.Sprintf("[GIN-debug] "+format, values...))
		}
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			global.LOG.Debug("[GIN-debug] print route",
				slog.String("method", httpMethod),
				slog.String("path", absolutePath),
				slog.String("name", handlerName),
				slog.Int("len", nuHandlers),
			)
		}
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	apis.SetupTracing()
	engine := gin.New()
	engine.RouterGroup.Use(otelgin.Middleware("")) // If the primary server name is not known, the default req.Host is used
	engine.RouterGroup.Use(Logger(global.LOG))
	engine.RouterGroup.Use(Recovery(global.LOG))
	engine.NoRoute()
	engine.NoMethod()

	// RouterPrefix: /file/
	engine.Handle(http.MethodGet, "/file/objects", file.ListHandler)
	engine.Handle(http.MethodPost, "/file/objects", file.CreateHandler)
	engine.Handle(http.MethodGet, "/file/objects/:id", file.DownloadHandler)
	engine.Handle(http.MethodDelete, "/file/objects/:id", file.DeleteHandler)

	return engine
}

// https://github.com/gin-contrib/zap/blob/173fe6c2ee6cd30c8891d5cf956e35b32e01dd0c/zap.go#L57
func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// log request info
		fields := []any{
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Duration("latency", latency),
		}

		// log request body
		// var buf bytes.Buffer
		// tee := io.TeeReader(c.Request.Body, &buf)
		// body, _ := io.ReadAll(tee)
		// c.Request.Body = io.NopCloser(&buf)
		// fields = append(fields, slog.String("body", string(body)))

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.ErrorContext(c.Request.Context(), e, fields...)
			}
		} else {
			logger.InfoContext(c.Request.Context(), "", fields...)
		}
	}
}

// https://github.com/gin-contrib/zap/blob/173fe6c2ee6cd30c8891d5cf956e35b32e01dd0c/sloggo#L145
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.ErrorContext(c.Request.Context(), c.Request.URL.Path,
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				logger.ErrorContext(c.Request.Context(), "[Recovery from panic]",
					slog.Any("error", err),
					slog.String("request", string(httpRequest)),
					slog.String("stack", string(debug.Stack())),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
