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
	}

	admin := v1.Group("/admin")
	{
		// admin 登录, 服务器采用 https, 因此密码传输在前端就不再额外处理了
		admin.POST("/login", handler.AdminLoginHandler(adminConf))
		// 新增paper
		admin.POST("/add_paper", handler.AddNewPaperHandler)
		// 修改paper
		admin.POST("/edit_paper", handler.EditPaperHandler)
	}
}
