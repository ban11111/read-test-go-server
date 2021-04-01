package handler

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"read-test-server/model"
	"read-test-server/service"

	"github.com/gin-gonic/gin"
	"read-test-server/common"
)

var ErrUploadFailed = errors.New("audio upload failed, please refresh")

// 上传音频接口
func UploadHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("record")
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	req, err := parseUploadParam(c)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	fileUrl, err := service.SaveFile(fileHeader, req)
	if err != nil {
		common.Log.Error("saveFile failed", zap.Error(err))
		common.RenderFail(c, ErrUploadFailed)
		return
	}

	fmt.Println("fileUrl:", fileUrl)
	// todo: 保存数据

	common.RenderSuccess(c)
}

func parseUploadParam(c *gin.Context) (*model.UploadReq, error) {
	var req model.UploadReq
	if paperName := c.PostForm("paper_name"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param paper_name is empty")
	}
	if paperName := c.PostForm("paper_version"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param paper_version is empty")
	}
	if paperName := c.PostForm("uid"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param uid is empty")
	}
	if paperName := c.PostForm("word_index"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param word_index is empty")
	}
	if paperName := c.PostForm("word"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param word is empty")
	}
	if paperName := c.PostForm("duration"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param duration is empty")
	}
	if paperName := c.PostForm("file_ext"); paperName != "" {
		req.PaperName = paperName
	} else {
		return nil, errors.New("param file_ext is empty")
	}
	return &req, nil
}

func SignInHandler(c *gin.Context) {
	// todo
}

func SignUpHandler(c *gin.Context) {
	// todo
	if err := service.SignUp(); err != nil {
		common.RenderFail(c, err)
		return
	}
}
