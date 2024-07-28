package apis

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func SetupTracing() {
	if tracer == nil {
		tracer = otel.GetTracerProvider().Tracer("github.com/whoisnian/k8s-example/src/user/pkg/apis")
	}
}

func JSON(c *gin.Context, code int, obj any) {
	_, span := tracer.Start(c.Request.Context(), "gin.render.json")
	defer span.End()
	c.Render(code, render.JSON{Data: obj})
}

func JSONMessage(c *gin.Context, code int, msg Message) {
	JSON(c, code, MessageResponse{Message: msg})
}

func AbortWithJSON(c *gin.Context, code int, obj any) {
	c.Abort()
	JSON(c, code, obj)
}

func AbortWithJSONMessage(c *gin.Context, code int, msg Message) {
	c.Abort()
	JSON(c, code, MessageResponse{Message: msg})
}

func String(c *gin.Context, code int, format string, values ...any) {
	_, span := tracer.Start(c.Request.Context(), "gin.render.string")
	defer span.End()
	c.Render(code, render.String{Format: format, Data: values})
}

func Redirect(c *gin.Context, code int, location string) {
	_, span := tracer.Start(c.Request.Context(), "gin.render.redirect")
	defer span.End()
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

func DataFromReader(c *gin.Context, code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	_, span := tracer.Start(c.Request.Context(), "gin.render.reader")
	defer span.End()
	c.Render(code, render.Reader{
		Headers:       extraHeaders,
		ContentType:   contentType,
		ContentLength: contentLength,
		Reader:        reader,
	})
}
