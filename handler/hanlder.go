package handler

import (
	"errors"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
	"read-test-server/common"
)

var ErrUploadFailed = errors.New("audio upload failed, please refresh")

// 保存文件 或者 上传 s3?  s3 国内访问太拉胯了, 先直接保存吧
func saveFile(fileHeader *multipart.FileHeader) error {
	fileName := fileHeader.Filename
	formFile, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer formFile.Close()
	sysFile, err := os.Create(common.AudioUploadRoot + fileName)
	if err != nil {
		return err
	}
	defer sysFile.Close()
	_, err = io.Copy(sysFile, formFile)
	if err != nil {
		return err
	}
	return nil
}

// 上传音频接口
func UploadHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("record")
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	err = saveFile(fileHeader)
	if err != nil {
		common.Log.Error("saveFile failed", zap.Error(err))
		common.RenderFail(c, ErrUploadFailed)
		return
	}

	common.RenderSuccess(c)
}
