package model

import (
	"errors"
	"fmt"
	"path"
	"regexp"
)

// 用户作答上传参数
type UploadReq struct {
	PaperName    string `json:"paper_name"`    // 用来生成文件路径和文件名的
	PaperVersion string `json:"paper_version"` // 用来生成文件路径和文件名的
	FileExt      string `json:"file_ext"`      // 用来生成文件路径和文件名的
	PaperId      uint   `json:"paper_id"`
	Uid          uint   `json:"uid"`
	WordIndex    int    `json:"word_index"`
	Word         string `json:"word"`
	Translation  string `json:"translation"` // 可以为空
	Duration     int    `json:"duration"`    // 单位: 毫秒
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

func (req *SignUpReq) ParamCheck() error {
	if req.Email == "" {
		return errors.New("please input your Email")
	}
	if req.Name == "" {
		return errors.New("please input your Name")
	}
	if req.ChineseClass == "" {
		return errors.New("please input your ChineseClass")
	}
	if req.HksLevel == "" {
		return errors.New("please input your HksLevel")
	}
	if req.EthnicBackground == "" {
		return errors.New("please input your EthnicBackground")
	}
	return nil
}

// 登录请求参数
type SignInReq struct {
	Email string `json:"email"`
}

func (req *SignInReq) ParamCheck() error {
	if req.Email == "" {
		return errors.New("please input your Email")
	}
	return nil
}

// 登录返回参数
type SignInResp struct {
	UserNotExist bool `json:"user_not_exist"` // 如果用户不存在, 则为 true
}

type AdminLogin struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
