package router

import (
	"github.com/gin-gonic/gin"
	"read-test-server/common"
	"read-test-server/handler"
)

func RegisterRouter(r *gin.Engine) {

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
		// todo 新增或编辑paper
		admin.POST("/add_or_update_paper")
		// todo
		admin.POST("/add_or_update_paper")
	}
}
