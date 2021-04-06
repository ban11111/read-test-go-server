package service

import (
	"errors"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path"
	"read-test-server/common"
	"read-test-server/dao"
	"read-test-server/model"
)

// 保存文件 或者 上传 s3?  s3 国内访问太拉胯了, 先直接保存吧, 30G硬盘应该够了
// ps: 文件路径说明：$root / 试卷名 / 用户id / 文件名
// 文件名说明: 试卷名-试卷版本-用户id-题目顺序-字词-耗时.mp3
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
func SignUp(req *model.SignUpReq) (err error) {
	err = dao.CreateUser(&model.User{
		Email:            req.Email,
		Name:             req.Name,
		ChineseClass:     req.ChineseClass,
		HksLevel:         req.HksLevel,
		EthnicBackground: req.EthnicBackground,
	})
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
