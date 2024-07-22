package user

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/user/global"
	"github.com/whoisnian/k8s-example/src/user/model"
	"github.com/whoisnian/k8s-example/src/user/pkg/apis"
	"go.uber.org/zap"
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
		c.AbortWithStatusJSON(http.StatusForbidden, apis.MessageResponse{Message: apis.MsgDisableRegistration})
		return
	}

	params := SignUpParams{}
	if err := c.BindJSON(&params); err != nil {
		global.LOG.Error("BindJSON SignUpParams", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgInvalidParams})
		return
	}

	var exists int64
	if err := global.DB.WithContext(c.Request.Context()).Model(&model.User{}).Where("email = ?", params.Email).Select("1").Find(&exists).Error; err != nil {
		global.LOG.Error("db find user", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}
	if exists == 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgInvalidEmail})
		return
	}

	// https://github.com/bcrypt-ruby/bcrypt-ruby/blob/master/lib/bcrypt/engine.rb
	// on my 2017 laptop with Intel Core i5-7300U
	// BCrypt::Engine.calibrate(55ms)  => 10
	// BCrypt::Engine.calibrate(220ms) => 12
	digest, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		global.LOG.Error("bcrypt generate", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	user := model.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: string(digest),
	}
	if err = global.DB.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		global.LOG.Error("db create user", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	c.JSON(http.StatusOK, user.AsJson())
}

type SignInParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func SignInHandler(c *gin.Context) {
	params := SignInParams{}
	if err := c.BindJSON(&params); err != nil {
		global.LOG.Error("BindJSON SignInParams", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgInvalidParams})
		return
	}

	var user model.User
	err := global.DB.WithContext(c.Request.Context()).First(&user, "email = ? AND deleted_at IS NULL", params.Email).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusForbidden, apis.MessageResponse{Message: apis.MsgEmailOrPasswordError})
		return
	} else if err != nil {
		global.LOG.Error("db find user", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		c.AbortWithStatusJSON(http.StatusForbidden, apis.MessageResponse{Message: apis.MsgEmailOrPasswordError})
		return
	} else if err != nil {
		global.LOG.Error("compare password", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err = session.Save(); err != nil {
		global.LOG.Error("save session", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	c.JSON(http.StatusOK, user.AsJson())
}

func LogoutHandler(c *gin.Context) {
	next := c.Query("next")
	if next == "" {
		next = "/"
	}

	session := sessions.Default(c)
	session.Delete("user_id")
	if err := session.Save(); err != nil {
		global.LOG.Error("save session", zap.Error(err))
	}
	c.Redirect(http.StatusTemporaryRedirect, next)
}

func InfoHandler(c *gin.Context) {
	// maybe unexpected log format:
	// https://github.com/gin-contrib/sessions/blob/4814ef52395a0762cd27afc049b3b38e56a28abe/sessions.go#L134
	id, ok := sessions.Default(c).Get("user_id").(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, apis.MessageResponse{Message: apis.MsgInvalidCookie})
		return
	}

	var user model.User
	err := global.DB.WithContext(c.Request.Context()).First(&user, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, apis.MessageResponse{Message: apis.MsgInvalidCookie})
		return
	} else if err != nil {
		global.LOG.Error("db find user", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	c.JSON(http.StatusOK, user.AsJson())
}

func InternalInfoHandler(c *gin.Context) {
	// maybe unexpected log format:
	// https://github.com/gin-contrib/sessions/blob/4814ef52395a0762cd27afc049b3b38e56a28abe/sessions.go#L134
	id, ok := sessions.Default(c).Get("user_id").(int64)
	if !ok {
		c.JSON(http.StatusOK, apis.NewInternalResponse[*model.User](apis.CodeInvalidCookie, apis.MsgInvalidCookie, nil))
		return
	}

	var user model.User
	err := global.DB.WithContext(c.Request.Context()).First(&user, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, apis.NewInternalResponse[*model.User](apis.CodeInvalidCookie, apis.MsgInvalidCookie, nil))
		return
	} else if err != nil {
		global.LOG.Error("db find user", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.NewInternalResponse[*model.User](apis.CodeInternalError, apis.MsgInternalError, nil))
		return
	}

	c.JSON(http.StatusOK, apis.NewInternalResponse(apis.CodeSuccess, apis.MsgSuccess, &user))
}
