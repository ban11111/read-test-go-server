package model

import (
	"fmt"
	"path"
	"regexp"
)

// 用户作答上传参数
type UploadReq struct {
	PaperName    string `json:"paper_name"`
	PaperVersion string `json:"paper_version"`
	Uid          uint   `json:"uid"`
	WordIndex    int    `json:"word_index"`
	Word         string `json:"word"`
	Duration     string `json:"duration"`
	FileExt      string `json:"file_ext"`
}

// 试卷名-试卷版本-用户id-题目顺序-字词-耗时
func (req *UploadReq) GetFileName() string {
	reg := regexp.MustCompile("[ \t]+")
	paperName := reg.ReplaceAllString(req.PaperName, "_")
	return fmt.Sprintf("%s-%s-%d-%d-%s-%s.%s", paperName, req.PaperVersion, req.Uid,
		req.WordIndex, req.Word, req.Duration, req.FileExt)
}

func (req *UploadReq) GetMiddlePath() string {
	reg := regexp.MustCompile("[ \t]+")
	paperName := reg.ReplaceAllString(req.PaperName, "_")
	return path.Join(paperName, fmt.Sprintf("%d", req.Uid))
}

// 注册请求参数
type SignUpReq struct {
	Email            string `json:"email"`
	Name             string `json:"name"`
	ChineseClass     string `json:"chinese_class"`
	HksLevel         string `json:"hks_level"`
	EthnicBackground string `json:"ethnic_background"`
}

// 登录请求参数
type SignInReq struct {
	Email string `json:"email"`
}

// 登录返回参数
type SignInResp struct {
	UserNotExist bool `json:"user_not_exist"` // 如果用户不存在, 则为 true
}
