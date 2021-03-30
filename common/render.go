package common

import "github.com/gin-gonic/gin"

func RenderSuccess(c *gin.Context, data ...interface{}) {
	if len(data) <= 0 {
		c.JSON(200, gin.H{"success": true})
	} else { // 只取一个
		c.JSON(200, gin.H{"success": true, "data": data[0]})
	}
}

func RenderFail(c *gin.Context, err error) {
	c.JSON(400, gin.H{"success": false, "info": err.Error()})
}
