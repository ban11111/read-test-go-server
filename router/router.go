package router

import (
	"github.com/gin-gonic/gin"
	"read-test-server/common"
	"read-test-server/handler"
)

func RegisterRouter(r *gin.Engine, adminConf *common.AdminConfig) {

	r.Static("file", common.AudioUploadRoot)
	v1 := r.Group("/api/v1")
	{
		// 上传录音文件,
		// ps: 因为需要时时保存进度方便从断点继续答题, 因此
		v1.POST("upload", handler.UploadHandler)
		// 登录
		v1.POST("sign_in", handler.SignInHandler)
		// 注册
		v1.POST("sign_up", handler.SignUpHandler)
		// 获取基本信息
		v1.POST("get_basic_info", handler.GetBasicInfoHandler)
	}

	admin := v1.Group("/admin")
	{
		// admin 登录, 服务器采用 https, 因此密码传输在前端就不再额外处理了
		admin.POST("/login", handler.AdminLoginHandler(adminConf))
	}
	adminAuth := v1.Group("/admin")
	adminAuth.Use(common.FakeTokenMiddleware(adminConf))
	{
		// 新增paper
		adminAuth.POST("/add_paper", handler.AddNewPaperHandler)
		// 修改paper
		adminAuth.POST("/edit_paper", handler.EditPaperHandler)
		// 修改paper
		adminAuth.POST("/query_papers", handler.QueryPapersHandler)
		// 查询用户列表
		adminAuth.POST("/query_users", handler.QueryUsersHandler)
		// 查询用户列表
		adminAuth.POST("/delete_user", handler.DeleteUsersHandler)
		// 查询用户做过的试卷
		adminAuth.POST("/query_user_paper")
		// 查询答题结果 (用户-试卷 维度)
		adminAuth.POST("/query_answers", handler.QueryAnswersHandler)
		// 删除答题结果 (用户-试卷 维度)
		adminAuth.POST("/clear_answers", handler.ClearAnswersHandler)
		// 更新 global setting
		adminAuth.POST("/query_settings", handler.QueryGlobalSettingsHandler)
		// 更新 global setting
		adminAuth.POST("/update_setting", handler.UpdateGlobalSettingHandler)
	}
}
