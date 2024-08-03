package user

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/user/global"
	"github.com/whoisnian/k8s-example/src/user/model"
	"github.com/whoisnian/k8s-example/src/user/pkg/apis"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignUpParams struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func SignUpHandler(c *gin.Context) {
	if global.CFG.DisableRegistration {
		apis.AbortWithJSONMessage(c, http.StatusForbidden, apis.MsgDisableRegistration)
		return
	}

	params := SignUpParams{}
	if err := c.BindJSON(&params); err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "BindJSON SignUpParams", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}

	var exists int64
	if err := global.DB.WithContext(c.Request.Context()).Model(&model.User{}).Where("email = ?", params.Email).Select("1").Find(&exists).Error; err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "db find user", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	if exists == 1 {
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidEmail)
		return
	}

	// https://github.com/bcrypt-ruby/bcrypt-ruby/blob/master/lib/bcrypt/engine.rb
	// on my 2017 laptop with Intel Core i5-7300U
	// BCrypt::Engine.calibrate(55ms)  => 10
	// BCrypt::Engine.calibrate(220ms) => 12
	digest, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "bcrypt generate", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	user := model.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: string(digest),
	}
	if err = global.DB.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "db create user", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	apis.JSON(c, http.StatusOK, user.AsJson())
}

type SignInParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignInHandler(c *gin.Context) {
	params := SignInParams{}
	if err := c.BindJSON(&params); err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "BindJSON SignInParams", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}

	var user model.User
	err := global.DB.WithContext(c.Request.Context()).First(&user, "email = ? AND deleted_at IS NULL", params.Email).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		apis.AbortWithJSONMessage(c, http.StatusForbidden, apis.MsgEmailOrPasswordError)
		return
	} else if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "db find user", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		apis.AbortWithJSONMessage(c, http.StatusForbidden, apis.MsgEmailOrPasswordError)
		return
	} else if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "compare password", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err = session.Save(); err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "save session", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	apis.JSON(c, http.StatusOK, user.AsJson())
}

func LogoutHandler(c *gin.Context) {
	next := c.Query("next")
	if next == "" {
		next = "/"
	} else if u, err := url.Parse(next); err != nil {
		next = "/"
		global.LOG.WarnContext(c.Request.Context(), "logout parse next", slog.Any("error", err))
	} else if u.Scheme != "" || u.Host != "" {
		next = "/"
		global.LOG.WarnContext(c.Request.Context(), "logout open-redirect detected", slog.String("nextScheme", u.Scheme), slog.String("nextHost", u.Host))
	}

	session := sessions.Default(c)
	session.Delete("user_id")
	if err := session.Save(); err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "save session", slog.Any("error", err))
	}
	apis.Redirect(c, http.StatusTemporaryRedirect, next)
}

func InfoHandler(c *gin.Context) {
	// maybe unexpected log format:
	// https://github.com/gin-contrib/sessions/blob/4814ef52395a0762cd27afc049b3b38e56a28abe/sessions.go#L134
	id, ok := sessions.Default(c).Get("user_id").(int64)
	if !ok {
		apis.AbortWithJSONMessage(c, http.StatusUnauthorized, apis.MsgInvalidCookie)
		return
	}

	var user model.User
	err := global.DB.WithContext(c.Request.Context()).First(&user, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		apis.AbortWithJSONMessage(c, http.StatusUnauthorized, apis.MsgInvalidCookie)
		return
	} else if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "db find user", slog.Any("error", err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	apis.JSON(c, http.StatusOK, user.AsJson())
}

func InternalInfoHandler(c *gin.Context) {
	// maybe unexpected log format:
	// https://github.com/gin-contrib/sessions/blob/4814ef52395a0762cd27afc049b3b38e56a28abe/sessions.go#L134
	id, ok := sessions.Default(c).Get("user_id").(int64)
	if !ok {
		apis.JSON(c, http.StatusOK, apis.NewInternalResponse[*model.User](apis.CodeInvalidCookie, apis.MsgInvalidCookie, nil))
		return
	}

	var user model.User
	err := global.DB.WithContext(c.Request.Context()).First(&user, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		apis.JSON(c, http.StatusOK, apis.NewInternalResponse[*model.User](apis.CodeInvalidCookie, apis.MsgInvalidCookie, nil))
		return
	} else if err != nil {
		global.LOG.ErrorContext(c.Request.Context(), "db find user", slog.Any("error", err))
		apis.AbortWithJSON(c, http.StatusInternalServerError, apis.NewInternalResponse[*model.User](apis.CodeInternalError, apis.MsgInternalError, nil))
		return
	}

	apis.JSON(c, http.StatusOK, apis.NewInternalResponse(apis.CodeSuccess, apis.MsgSuccess, &user))
}
