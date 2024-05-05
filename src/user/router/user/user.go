package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func SigninHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
