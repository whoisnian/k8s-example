package file

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/file/global"
	"github.com/whoisnian/k8s-example/src/file/model"
	"github.com/whoisnian/k8s-example/src/file/pkg/apis"
	"github.com/whoisnian/k8s-example/src/file/service/svcuser"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ListHandler(c *gin.Context) {
	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	var files []model.File
	err = global.DB.WithContext(c.Request.Context()).Where("user_id = ? AND deleted_at IS NULL", user.ID).Order("id desc").Find(&files).Error
	if err != nil {
		global.LOG.Error("db find files", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}
	results := make([]model.FileJson, len(files))
	for i := range files {
		results[i] = files[i].AsJson()
	}
	c.JSON(http.StatusOK, results)
}

func CreateHandler(c *gin.Context) {
	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		global.LOG.Error("get multipart reader", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgMultipartReaderError})
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			global.LOG.Error("read multipart reader part", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, apis.MessageResponse{Message: apis.MsgMultipartReaderError})
			return
		}

		var sizes []int64
		if part.FormName() == "fileSize" {
			err = json.NewDecoder(part).Decode(&sizes)
			if err != nil {
				global.LOG.Error("read multipart reader fileSize", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgInvalidParams})
				return
			}
		} else if part.FormName() == "fileList" {
			global.LOG.Debug("multipart", zap.Any("part", part))
			file := model.File{
				UserID:     user.ID,
				Name:       part.FileName(),
				BucketName: global.CFG.StorageBucket,
			}
			err = global.DB.WithContext(c.Request.Context()).Create(&file).Error
			if err != nil {
				global.LOG.Error("db create file", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
				return
			}
			file.ObjectName = genObjectName(file.ID)

			hasher := sha256.New()
			var size int64 = -1
			if len(sizes) > 0 {
				size = sizes[0]
				sizes = sizes[1:]
			}
			file.Size, err = global.FS.CreateFile(file.BucketName, file.ObjectName, io.TeeReader(part, hasher), size)
			if err != nil {
				global.LOG.Error("fs create file", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
				return
			}
			file.Digest = hex.EncodeToString(hasher.Sum(nil))

			err = global.DB.WithContext(c.Request.Context()).Save(&file).Error
			if err != nil {
				global.LOG.Error("db update file", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
				return
			}
		}
	}
	c.JSON(http.StatusOK, apis.MessageResponse{Message: apis.MsgSuccess})
}

func DownloadHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgInvalidParams})
		return
	}

	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	var file model.File
	err = global.DB.WithContext(c.Request.Context()).First(&file, "user_id = ? AND id = ? AND deleted_at IS NULL", user.ID, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, apis.MessageResponse{Message: apis.MsgFileNotFound})
		return
	} else if err != nil {
		global.LOG.Error("db find file", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}
	if file.ObjectName == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, apis.MessageResponse{Message: apis.MsgFileIncomplete})
		return
	}

	irc, err := global.FS.OpenFile(file.BucketName, file.ObjectName)
	if err != nil {
		global.LOG.Error("fs open file", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}
	defer irc.Close()

	c.DataFromReader(http.StatusOK, file.Size, "application/octet-stream", irc, nil)
}

func DeleteHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apis.MessageResponse{Message: apis.MsgInvalidParams})
		return
	}

	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	var file model.File
	err = global.DB.WithContext(c.Request.Context()).First(&file, "user_id = ? AND id = ? AND deleted_at IS NULL", user.ID, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, apis.MessageResponse{Message: apis.MsgSuccess})
		return
	} else if err != nil {
		global.LOG.Error("db find file", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}

	file.DeletedAt.Scan(time.Now())
	err = global.DB.WithContext(c.Request.Context()).Save(&file).Error
	if err != nil {
		global.LOG.Error("db delete file", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, apis.MessageResponse{Message: apis.MsgInternalError})
		return
	}
	c.JSON(http.StatusOK, apis.MessageResponse{Message: apis.MsgSuccess})
}

func genObjectName(id int64) string {
	end := id & 31 // id % 32
	id = id >> 5   // id / 32
	mid := id & 31 // id % 32
	id = id >> 5   // id / 32

	buf := make([]byte, 0, 16)
	buf = strconv.AppendInt(buf, id, 32)
	buf = append(buf, '/')
	buf = strconv.AppendInt(buf, mid, 32)
	buf = append(buf, '/')
	buf = strconv.AppendInt(buf, end, 32)
	return string(buf)
}
