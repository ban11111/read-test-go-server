package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"read-test-server/model"
	"strconv"
)

func parseUploadParam(c *gin.Context) (*model.UploadReq, error) {
	var req model.UploadReq
	if req.PaperName = c.PostForm("paper_name"); req.PaperName == "" {
		return nil, errors.New("param paper_name is empty")
	}
	if req.PaperVersion = c.PostForm("paper_version"); req.PaperVersion == "" {
		return nil, errors.New("param paper_version is empty")
	}
	if req.FileExt = c.PostForm("file_ext"); req.FileExt == "" {
		return nil, errors.New("param file_ext is empty")
	}
	if paperId := c.PostForm("paper_id"); paperId != "" {
		parseUint, err := strconv.ParseUint(paperId, 10, 64)
		if err != nil {
			return nil, errors.New("param paper_id is not a valid number")
		}
		if parseUint <= 0 {
			return nil, errors.New("param paper_id is empty")
		}
		req.PaperId = uint(parseUint)
	} else {
		return nil, errors.New("param paper_id is empty")
	}
	if uid := c.PostForm("uid"); uid != "" {
		parseUint, err := strconv.ParseUint(uid, 10, 64)
		if err != nil {
			return nil, errors.New("param uid is not a valid number")
		}
		if parseUint <= 0 {
			return nil, errors.New("param uid is empty")
		}
		req.Uid = uint(parseUint)
	} else {
		return nil, errors.New("param uid is empty")
	}
	if wordIndex := c.PostForm("word_index"); wordIndex != "" {
		parseInt, err := strconv.ParseInt(wordIndex, 10, 64)
		if err != nil {
			return nil, errors.New("param word_index is not a valid number")
		}
		if parseInt < 0 {
			return nil, errors.New("param word_index is empty")
		}
		req.WordIndex = int(parseInt)
	} else {
		return nil, errors.New("param word_index is empty")
	}
	if req.Word = c.PostForm("word"); req.Word == "" {
		return nil, errors.New("param word is empty")
	}
	req.Translation = c.PostForm("translation")
	if duration := c.PostForm("duration"); duration != "" {
		parseInt, err := strconv.ParseInt(duration, 10, 64)
		if err != nil {
			return nil, errors.New("param duration is not a valid number")
		}
		if parseInt <= 0 {
			return nil, errors.New("param duration is empty")
		}
		req.Duration = int(parseInt)
	} else {
		return nil, errors.New("param duration is empty")
	}
	req.Device = c.PostForm("device")
	return &req, nil
}