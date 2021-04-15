package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
	var withFile = true
	fileHeader, err := c.FormFile("record")
	if errors.Is(err, http.ErrMissingFile) {
		withFile = false
	}
	if err != nil && withFile {
		common.RenderFail(c, err)
		return
	}
	req, err := parseUploadParam(c)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	audioUrl := ""
	if withFile {
		audioUrl, err = service.SaveFile(fileHeader, req)
		if err != nil {
			common.Log.Error("saveFile failed", zap.Error(err))
			common.RenderFail(c, ErrUploadFailed)
			return
		}
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
	user, err := service.SignUp(&req)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, map[string]interface{}{"user": user})
}

func GetBasicInfoHandler(c *gin.Context) {
	var req model.BasicInfoReq
	if err := c.BindJSON(&req); err != nil {
		common.RenderFail(c, ErrParamInvalid)
		return
	}

	resp, err := service.GetBasicInfo(&req)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, &resp)
}

// admin 登录接口
func AdminLoginHandler(adminConf *common.AdminConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req model.AdminLogin
		if err := c.BindJSON(&req); err != nil {
			common.RenderFail(c, ErrParamInvalid)
			return
		}
		if req.Username != adminConf.Username {
			common.Log.Error("AdminLoginHandler", zap.String("req.Username", req.Username), zap.String("conf.Username", adminConf.Username))
			common.RenderFail(c, ErrWrongUserName)
			return
		}
		if !common.MatchPass(req.Password, adminConf.EncodedPassword) {
			common.RenderFail(c, ErrWrongPassword)
			return
		}
		// 偷懒了, token 直接放加密后的密码得了
		common.RenderSuccess(c, gin.H{"token": adminConf.EncodedPassword})
	}
}

// ================= todo 试卷相关接口 ===================

// 新增试卷
func AddNewPaperHandler(c *gin.Context) {

}

// 修改试卷
// ps: 每次修改其实是重新创建一个试卷, 并且版本号更新
func EditPaperHandler(c *gin.Context) {
	var req model.Paper
	if err := c.BindJSON(&req); err != nil {
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	if err := service.EditPaper(&req); err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c)
}

// 查询试卷列表, 直接 按id倒序 全部查出来扔前端就行, 暂时不考虑分页
func QueryPapersHandler(c *gin.Context) {
	papers, activePaper, err := service.QueryPapers()
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, gin.H{"papers": papers, "active_paper": activePaper})
}

func QueryUsersHandler(c *gin.Context) {
	users, err := service.QueryUsers()
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, users)
}

func DeleteUsersHandler(c *gin.Context) {
	var req model.DeleteUserReq
	if err := c.BindJSON(&req); err != nil {
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	err := service.DeleteUser(req.Uid)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c)
}

func QueryAnswersHandler(c *gin.Context) {
	var req model.GetAnswersReq
	if err := c.BindJSON(&req); err != nil {
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	answers, err := service.QueryAnswers(&req)
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, answers)
}

func ClearAnswersHandler(c *gin.Context) {
	var req model.GetAnswersReq
	if err := c.BindJSON(&req); err != nil {
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	if err := service.ClearAnswers(&req); err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c)
}

func QueryGlobalSettingsHandler(c *gin.Context) {
	settings, err := service.QueryGlobalSettings()
	if err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c, settings)
}

func UpdateGlobalSettingHandler(c *gin.Context) {
	var req map[string]interface{}

	if err := c.BindJSON(&req); err != nil {
		common.RenderFail(c, ErrParamInvalid)
		return
	}
	if err := service.UpdateGlobalSetting(req); err != nil {
		common.RenderFail(c, err)
		return
	}
	common.RenderSuccess(c)
}

// ================= todo 统计相关接口 ===================

func GetStatistics(c *gin.Context) {

}
