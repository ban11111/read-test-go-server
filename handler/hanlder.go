package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"read-test-server/common"
	"read-test-server/model"
	"read-test-server/service"
)

// 这边定义返回给用户的报错信息
var ErrUploadFailed = errors.New("uploading audio failed, please refresh or try again later")
var ErrSaveAnswerFailed = errors.New("saving answer failed, please refresh or try again later")
var ErrParamInvalid = errors.New("param invalid, please contact web manager")
var ErrWrongUserName = errors.New("wrong user_name")
var ErrWrongPassword = errors.New("wrong password")

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
	audioUrl, err := service.SaveFile(fileHeader, req)
	if err != nil {
		common.Log.Error("saveFile failed", zap.Error(err))
		common.RenderFail(c, ErrUploadFailed)
		return
	}

	common.Log.Debug("audioUrl:", zap.String("url", audioUrl))
	if err = service.SaveAnswer(req, audioUrl); err != nil {
		common.RenderFail(c, ErrSaveAnswerFailed)
		return
	}

	common.RenderSuccess(c)
}

func SignInHandler(c *gin.Context) {
	var req model.SignInReq
	if err := c.BindJSON(&req); err != nil {
		common.Log.Error("SignInHandler.BindJSON()", zap.Error(err))
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	if err := req.ParamCheck(); err != nil {
		common.RenderFail(c, err)
		return
	}
	userNotExist, user, err := service.SignIn(req.Email)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, map[string]interface{}{"user_not_exist": userNotExist, "user": user})
}

func SignUpHandler(c *gin.Context) {
	var req model.SignUpReq
	if err := c.BindJSON(&req); err != nil {
		common.Log.Error("SignUpHandler.BindJSON()", zap.Error(err))
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	if err := req.ParamCheck(); err != nil {
		common.RenderFail(c, err)
		return
	}
	if err := service.SignUp(&req); err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c)
}

// admin 登录接口
func AdminLoginHandler(adminConf *common.AdminConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req model.AdminLogin
		if err := c.BindJSON(&req); err != nil {
			common.RenderFail(c, ErrParamInvalid)
			return
		}
		if req.UserName != adminConf.UserName {
			common.RenderFail(c, ErrWrongUserName)
			return
		}
		if !common.MatchPass(req.Password, adminConf.EncodedPassword) {
			common.RenderFail(c, ErrWrongPassword)
		}
		common.RenderSuccess(c)
	}
}

// ================= todo 试卷相关接口 ===================

// 新增试卷
func AddNewPaperHandler(c *gin.Context) {

}

// 修改试卷
// ps: 每次修改其实是重新创建一个试卷, 并且版本号更新
func EditPaperHandler(c *gin.Context) {

}

// 查询试卷列表, 直接 按id倒序 全部查出来扔前端就行, 暂时不考虑分页
func QueryPapersHandler(c *gin.Context) {

}

// ================= todo 统计相关接口 ===================

func GetStatistics(c *gin.Context) {

}
