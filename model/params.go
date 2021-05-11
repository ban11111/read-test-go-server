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
	Device       string `json:"device"`      // 设备类型
}

// 试卷名-试卷版本-用户id-题目顺序-字词
func (req *UploadReq) GetFileName() string {
	reg := regexp.MustCompile("[ \t]+")
	paperName := reg.ReplaceAllString(req.PaperName, "_")
	return fmt.Sprintf("%s-%s-%d-%d-%s.%s", paperName, req.PaperVersion, req.Uid,
		req.WordIndex, req.Word, req.FileExt)
}

func (req *UploadReq) GetMiddlePath() string {
	reg := regexp.MustCompile("[ \t]+")
	paperName := reg.ReplaceAllString(req.PaperName, "_")
	return path.Join(paperName, fmt.Sprintf("%d", req.Uid))
}

// 注册请求参数
type SignUpReq struct {
	Email                  string `json:"email"`
	Name                   string `json:"name"`
	ChineseClass           string `json:"chinese_class"`
	HksLevel               string `json:"hks_level"`
	EthnicBackground       string `json:"ethnic_background"`
	HasChineseAcquaintance bool   `json:"has_chinese_acquaintance"`
	AcquaintanceDetail     string `json:"acquaintance_detail"`
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
	Username string `json:"username"`
	Password string `json:"password"`
}

type BasicInfoReq struct {
	Uid uint `json:"uid"`
}

type BasicInfoResp struct {
	CurrentPaper  *Paper                 `json:"current_paper"`
	GlobalSetting map[string]interface{} `json:"global_setting"`
	ProgressIndex int                    `json:"progress_index"`
}

type GetAnswersReq struct {
	Uid          uint  `json:"uid"`
	PaperId      uint  `json:"paper_id"`
	PaperVersion int16 `json:"paper_version"`
}

type DeleteUserReq struct {
	Uid uint `json:"uid"`
}

type AddPaperReq struct {
	Name     string `json:"name"`
	Words    string `json:"words"`
	Interval int    `json:"interval"`
}

func (req *AddPaperReq) ParamCheck() error {
	if req.Name == "" {
		return errors.New("please input Paper Name")
	}
	if req.Interval <= 0 {
		return errors.New("please input Interval")
	}
	if req.Words == "" {
		return errors.New("please input Words")
	}
	return nil
}

type PublishPaperReq struct {
	Pid uint `json:"pid"`
}

func (req *PublishPaperReq) ParamCheck() error {
	if req.Pid <= 0 {
		return errors.New("empty pid")
	}
	return nil
}

type StatisticsResp struct {
	TotalUsers        int        `json:"total_users"`
	CurrentMonthUsers int        `json:"current_month_users"`
	TotalProgress     float32    `json:"total_progress"`
	TotalAnswers      int        `json:"total_answers"`
	Device            DeviceInfo `json:"device"`
	Chart             ChartInfo  `json:"chart"`
}

type DeviceInfo struct {
	Desktop int `json:"desktop"` // percentage
	Tablet  int `json:"tablet"`  // percentage
	Mobile  int `json:"mobile"`  // percentage
	Unknown int `json:"unknown"` // percentage
}

type ChartInfo struct {
	Daily   ChartCount `json:"daily"`   // 7 days
	Monthly ChartCount `json:"monthly"` // 7 months
}

type ChartCount struct {
	NewUser   []int    `json:"new_user"`
	NewAnswer []int    `json:"new_answer"`
	Labels    []string `json:"labels"`
}
