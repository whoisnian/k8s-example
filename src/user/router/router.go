package router

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/user/global"
	"github.com/whoisnian/k8s-example/src/user/pkg/apis"
	"github.com/whoisnian/k8s-example/src/user/router/user"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Setup() *gin.Engine {
	if global.CFG.Debug {
		gin.SetMode(gin.DebugMode)
		gin.DebugPrintFunc = func(format string, values ...interface{}) {
			global.LOG.Debug(fmt.Sprintf("[GIN-debug] "+format, values...))
		}
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			global.LOG.Debug("[GIN-debug] print route",
				zap.String("method", httpMethod),
				zap.String("path", absolutePath),
				zap.String("name", handlerName),
				zap.Int("len", nuHandlers),
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
	engine.RouterGroup.Use(sessions.Sessions("_app_session", global.RSS))
	engine.NoRoute()
	engine.NoMethod()

	// RouterPrefix: /user/
	engine.Handle(http.MethodPost, "/user/signup", user.SignUpHandler)
	engine.Handle(http.MethodPost, "/user/signin", user.SignInHandler)
	engine.Handle(http.MethodGet, "/user/logout", user.LogoutHandler)
	engine.Handle(http.MethodGet, "/user/info", user.InfoHandler)

	// RouterPrefix: /internal/user/
	engine.Handle(http.MethodGet, "/internal/user/info", user.InternalInfoHandler)

	return engine
}

// https://github.com/gin-contrib/zap/blob/173fe6c2ee6cd30c8891d5cf956e35b32e01dd0c/zap.go#L57
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// log request info
		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		// log opentelemetry trace
		if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}
		if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
			fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
			fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
		}

		// log request body
		// var buf bytes.Buffer
		// tee := io.TeeReader(c.Request.Body, &buf)
		// body, _ := io.ReadAll(tee)
		// c.Request.Body = io.NopCloser(&buf)
		// fields = append(fields, zap.String("body", string(body)))

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info("", fields...)
		}
	}
}

// https://github.com/gin-contrib/zap/blob/173fe6c2ee6cd30c8891d5cf956e35b32e01dd0c/zap.go#L145
func Recovery(logger *zap.Logger) gin.HandlerFunc {
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
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				logger.Error("[Recovery from panic]",
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
