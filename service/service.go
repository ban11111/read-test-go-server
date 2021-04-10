package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path"
	"read-test-server/common"
	"read-test-server/dao"
	"read-test-server/model"
	"strconv"
)

// 保存文件 或者 上传 s3?  s3 国内访问太拉胯了, 先直接保存吧, 30G硬盘应该够了
// ps: 文件路径说明：$root / 试卷名 / 用户id / 文件名
// 文件名说明: 试卷名-试卷版本-用户id-题目顺序-字词.mp3
func SaveFile(fileHeader *multipart.FileHeader, req *model.UploadReq) (fileUrl string, err error) {
	//fileName := fileHeader.Filename
	formFile, err := fileHeader.Open()
	if err != nil {
		return
	}
	defer formFile.Close()
	fileName := req.GetFileName()
	fileUrl = path.Join(req.GetMiddlePath(), fileName)
	filePath := path.Join(common.AudioUploadRoot, fileUrl)
	sysFile, err := os.Create(filePath)
	if err != nil {
		dir := path.Join(common.AudioUploadRoot, req.GetMiddlePath())
		_, err = os.Stat(dir)
		if err != nil {
			if err = os.MkdirAll(dir, 0777); err != nil {
				return
			}
			if sysFile, err = os.Create(filePath); err != nil {
				return
			}
		}
	}
	defer sysFile.Close()
	_, err = io.Copy(sysFile, formFile)
	if err != nil {
		return
	}
	return
}

// 保存每一题的答题结果
func SaveAnswer(req *model.UploadReq, audioUrl string) error {
	return dao.CreateAnswer(&model.Answer{
		PaperId:     req.PaperId,
		Uid:         req.Uid,
		WordIndex:   req.WordIndex,
		Word:        req.Word,
		AudioUrl:    audioUrl,
		Translation: req.Translation,
		Duration:    req.Duration,
	})
}

// 注册
func SignUp(req *model.SignUpReq) (user *model.User, err error) {
	user = &model.User{
		Email:            req.Email,
		Name:             req.Name,
		ChineseClass:     req.ChineseClass,
		HksLevel:         req.HksLevel,
		EthnicBackground: req.EthnicBackground,
	}
	err = dao.CreateUser(user)
	return
}

// 登录
func SignIn(email string) (needSignUp bool, user *model.User, err error) {
	user, err = dao.QueryUserByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil, nil
	}
	return
}

func GetBasicInfo(req *model.BasicInfoReq) (resp *model.BasicInfoResp, err error) {
	resp = &model.BasicInfoResp{
		GlobalSetting: make(map[string]interface{}),
	}
	if resp.CurrentPaper, err = dao.QueryCurrentPaper(); err != nil {
		return
	}
	settings, err := dao.QueryGlobalSettings()
	for _, set := range settings {
		parseInt, parseErr := strconv.ParseInt(set.Value, 10, 32)
		if parseErr != nil {
			err = parseErr
			return
		}
		resp.GlobalSetting[set.Key] = parseInt
	}
	progress, err := dao.QueryAnswerProgress(req.Uid, resp.CurrentPaper.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound){
		return
	}
	if progress != nil {
		resp.ProgressIndex = progress.WordIndex
	}
	return
}

func QueryPapers() (papers []*model.Paper, activePaper *model.Paper, err error) {
	papers, err = dao.QueryPapers()
	if err != nil {
		return
	}
	for _, p := range papers {
		if p.Inuse {
			activePaper = p
			return
		}
	}
	err = errors.New("no active paper in use, please contact zebreay! ")
	return
}

func EditPaper(paper *model.Paper) error {
	previousPaper, err := dao.QueryPaperById(paper.Id)
	if err != nil {
		return err
	}
	if previousPaper.Words == paper.Words && previousPaper.Interval == paper.Interval && previousPaper.Inuse == paper.Inuse {
		return errors.New("nothing to update")
	}
	if previousPaper.Words != paper.Words || previousPaper.Interval != paper.Interval {
		defer dao.CreatePaperSnapshot(&model.PaperSnapshot{Paper: model.Paper{
			Name:     previousPaper.Name,
			Version:  previousPaper.Version,
			Words:    previousPaper.Words,
			Interval: previousPaper.Interval,
			Inuse:    previousPaper.Inuse,
		}, PaperId: previousPaper.Id})
	}
	paper.Version++
	return dao.UpdatePaper(paper)
}

func QueryUsers() ([]*model.User, error) {
	return dao.QueryUsers()
}

func QueryAnswers(req *model.GetAnswersReq) ([]*model.Answer, error) {
	return dao.QueryAnswersByUidAndPaper(req.Uid, req.PaperId)
}

func QueryGlobalSettings() (map[string]interface{}, error) {
	settings, err := dao.QueryGlobalSettings()
	if err != nil {
		return nil, err
	}
	var result = make(map[string]interface{})
	for _, set := range settings {
		result[set.Key] = set.Value
	}
	return result, nil
}

func UpdateGlobalSetting(settings map[string]interface{}) error {
	existSettings, err := dao.QueryGlobalSettings()
	if err != nil {
		return err
	}
	for _, set := range existSettings {
		if newValue, ok := settings[set.Key]; ok {
			if err = dao.UpdateGlobalSetting(&model.GlobalSetting{
				Key:   set.Key,
				Value: fmt.Sprintf("%v", newValue),
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
