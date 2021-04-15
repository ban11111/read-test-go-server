package common

import "github.com/gin-gonic/gin"

// 懒得搞真的token了， 搞个假的算了
func FakeTokenMiddleware(adminConf *AdminConfig) func (c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token != adminConf.EncodedPassword {
			c.JSON(403, nil)
			c.Abort()
		}
	}
}