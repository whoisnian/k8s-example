package file

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/file/global"
	"github.com/whoisnian/k8s-example/src/file/model"
	"github.com/whoisnian/k8s-example/src/file/service/svcuser"
	"go.uber.org/zap"
)

func ListHandler(c *gin.Context) {
	info, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "svc err"})
		return
	}

	var files []model.File
	res := global.DB.WithContext(c.Request.Context()).Where("user_id = ? AND deleted_at IS NULL", info.ID).Order("id desc").Find(&files)
	if res.Error != nil {
		global.LOG.Error("db find files", zap.Error(res.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "db err"})
		return
	}
	c.JSON(http.StatusOK, files)
}

func CreateHandler(c *gin.Context) {
	info, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "svc err"})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		global.LOG.Error("get multipart reader", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "get multipart reader err"})
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			global.LOG.Error("read multipart reader", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"msg": "read multipart reader err"})
			return
		}

		var sizes []int64
		if part.FormName() == "fileSize" {
			err = json.NewDecoder(part).Decode(&sizes)
			if err != nil {
				global.LOG.Error("read multipart reader", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "read multipart reader err"})
				return
			}
		} else if part.FormName() == "fileList" {
			global.LOG.Debug("multipart", zap.Any("part", part))
			file := model.File{
				UserID:     info.ID,
				Name:       part.FileName(),
				BucketName: global.CFG.StorageBucket,
			}
			res := global.DB.WithContext(c.Request.Context()).Create(&file)
			if res.Error != nil {
				global.LOG.Error("db create file", zap.Error(res.Error))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "db err"})
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
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "fs err"})
				return
			}
			file.Digest = hex.EncodeToString(hasher.Sum(nil))

			res = global.DB.WithContext(c.Request.Context()).Save(&file)
			if res.Error != nil {
				global.LOG.Error("db update file", zap.Error(res.Error))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "db err"})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func DownloadHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	info, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "svc err"})
		return
	}

	var file model.File
	res := global.DB.WithContext(c.Request.Context()).First(&file, "user_id = ? AND id = ? AND deleted_at IS NULL", info.ID, id)
	if res.Error != nil {
		global.LOG.Error("db find file", zap.Error(res.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "db err"})
		return
	}

	if file.ObjectName == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"msg": "file is incomplete"})
		return
	}

	irc, err := global.FS.OpenFile(file.BucketName, file.ObjectName)
	if err != nil {
		global.LOG.Error("fs open file", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "fs err"})
		return
	}
	defer irc.Close()

	c.DataFromReader(http.StatusOK, file.Size, "application/octet-stream", irc, nil)
}

func DeleteHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	info, err := svcuser.UserInfo(c)
	if err != nil {
		global.LOG.Error("svc user info", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "svc err"})
		return
	}

	var file model.File
	res := global.DB.WithContext(c.Request.Context()).First(&file, "user_id = ? AND id = ? AND deleted_at IS NULL", info.ID, id)
	if res.Error != nil {
		global.LOG.Error("db find file", zap.Error(res.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "db err"})
		return
	}

	file.DeletedAt.Scan(time.Now())
	res = global.DB.WithContext(c.Request.Context()).Save(&file)
	if res.Error != nil {
		global.LOG.Error("db delete file", zap.Error(res.Error))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "db err"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
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
