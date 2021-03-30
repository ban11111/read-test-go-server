package router

import (
	"github.com/gin-gonic/gin"
	"read-test-server/handler"
)

func RegisterRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	// 上传录音文件,
	// ps: 因为需要时时保存进度方便从断点继续答题, 因此
	v1.POST("upload", handler.UploadHandler)
}
