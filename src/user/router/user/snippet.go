package user

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/whoisnian/k8s-example/src/user/global"
	"github.com/whoisnian/k8s-example/src/user/pkg/apis"
	"github.com/whoisnian/k8s-example/src/user/pkg/key"
)

func ShowSnippetHandler(c *gin.Context) {
	id, _ := sessions.Default(c).Get("user_id").(int64)
	content, err := global.RDB.Get(c.Request.Context(), key.RedisUserSnippet(id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		global.LOG.ErrorContext(c.Request.Context(), "redis get", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	dur, err := global.RDB.TTL(c.Request.Context(), key.RedisUserSnippet(id)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		global.LOG.ErrorContext(c.Request.Context(), "redis ttl", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	apis.JSON(c, http.StatusOK, map[string]any{
		"content":    content,
		"expiration": dur / time.Second,
	})
}

type UpdateSnippetParams struct {
	Content    string `json:"content" binding:"required"`
	Expiration string `json:"expiration" binding:"required"`
}

func UpdateSnippetHandler(c *gin.Context) {
	params := UpdateSnippetParams{}
	if err := c.BindJSON(&params); err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "BindJSON UpdateSnippetParams", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}
	seconds, err := strconv.ParseInt(params.Expiration, 10, 64)
	if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "strconv ParseInt", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}
	if len(params.Content) > 1024*1024*4 {
		apis.AbortWithJSONMessage(c, http.StatusRequestEntityTooLarge, apis.MsgInvalidParams)
		return
	}

	id, _ := sessions.Default(c).Get("user_id").(int64)
	err = global.RDB.Set(c.Request.Context(), key.RedisUserSnippet(id), params.Content, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "redis set", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	apis.JSON(c, http.StatusOK, apis.MsgSuccess)
}

func DeleteSnippetHandler(c *gin.Context) {
	id, _ := sessions.Default(c).Get("user_id").(int64)
	err := global.RDB.Del(c.Request.Context(), key.RedisUserSnippet(id)).Err()
	if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "redis del", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	apis.JSON(c, http.StatusOK, apis.MsgSuccess)
}
