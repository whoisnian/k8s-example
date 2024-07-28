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
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	var files []model.File
	err = global.DB.WithContext(c.Request.Context()).Where("user_id = ? AND deleted_at IS NULL", user.ID).Order("id desc").Find(&files).Error
	if err != nil {
		global.LOG.Error("db find files", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	results := make([]model.FileJson, len(files))
	for i := range files {
		results[i] = files[i].AsJson()
	}
	apis.JSON(c, http.StatusOK, results)
}

func CreateHandler(c *gin.Context) {
	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		global.LOG.Error("get multipart reader", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgMultipartReaderError)
		return
	}

	// https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#multipart-form-data
	// The order of parts must be the same as the order of fields in entry list.
	part, err := reader.NextPart()
	if err != nil {
		global.LOG.Error("read multipart reader part", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusUnprocessableEntity, apis.MsgMultipartReaderError)
		return
	} else if part.FormName() != "fileSize" {
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}

	var sizes []int64
	err = json.NewDecoder(part).Decode(&sizes)
	if err != nil {
		global.LOG.Error("read multipart reader fileSize", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			global.LOG.Error("read multipart reader part", zap.Error(err))
			apis.AbortWithJSONMessage(c, http.StatusUnprocessableEntity, apis.MsgMultipartReaderError)
			return
		} else if part.FormName() != "fileList" {
			apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
			return
		}

		global.LOG.Debug("multipart read", zap.Any("part", part))
		file := model.File{
			UserID:     -1,
			Name:       part.FileName(),
			BucketName: global.CFG.StorageBucket,
		}
		err = global.DB.WithContext(c.Request.Context()).Create(&file).Error
		if err != nil {
			global.LOG.Error("db create file", zap.Error(err))
			apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
			return
		}
		file.UserID = user.ID
		file.ObjectName = genObjectName(file.ID)

		hasher := sha256.New()
		var size int64 = -1
		if len(sizes) > 0 {
			size = sizes[0]
			sizes = sizes[1:]
		}
		file.Size, err = global.FS.CreateFile(c.Request.Context(), file.BucketName, file.ObjectName, io.TeeReader(part, hasher), size)
		if err != nil {
			global.LOG.Error("fs create file", zap.Error(err))
			apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
			return
		}
		file.Digest = hex.EncodeToString(hasher.Sum(nil))

		err = global.DB.WithContext(c.Request.Context()).Save(&file).Error
		if err != nil {
			global.LOG.Error("db update file", zap.Error(err))
			apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
			return
		}
	}
	apis.JSON(c, http.StatusOK, apis.MsgSuccess)
}

func DownloadHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}

	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	var file model.File
	err = global.DB.WithContext(c.Request.Context()).First(&file, "user_id = ? AND id = ? AND deleted_at IS NULL", user.ID, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		apis.AbortWithJSONMessage(c, http.StatusNotFound, apis.MsgFileNotFound)
		return
	} else if err != nil {
		global.LOG.Error("db find file", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	if file.ObjectName == "" {
		apis.AbortWithJSONMessage(c, http.StatusUnprocessableEntity, apis.MsgFileIncomplete)
		return
	}

	irc, err := global.FS.OpenFile(c.Request.Context(), file.BucketName, file.ObjectName)
	if err != nil {
		global.LOG.Error("fs open file", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	defer irc.Close()

	apis.DataFromReader(c, http.StatusOK, file.Size, "application/octet-stream", irc, nil)
}

func DeleteHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apis.AbortWithJSONMessage(c, http.StatusBadRequest, apis.MsgInvalidParams)
		return
	}

	user, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	var file model.File
	err = global.DB.WithContext(c.Request.Context()).First(&file, "user_id = ? AND id = ? AND deleted_at IS NULL", user.ID, id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		apis.JSONMessage(c, http.StatusOK, apis.MsgSuccess)
		return
	} else if err != nil {
		global.LOG.Error("db find file", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}

	file.DeletedAt.Scan(time.Now())
	err = global.DB.WithContext(c.Request.Context()).Save(&file).Error
	if err != nil {
		global.LOG.Error("db delete file", zap.Error(err))
		apis.AbortWithJSONMessage(c, http.StatusInternalServerError, apis.MsgInternalError)
		return
	}
	apis.JSONMessage(c, http.StatusOK, apis.MsgSuccess)
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
