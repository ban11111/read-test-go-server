package common

import (
	"github.com/dipperin/go-ms-toolkit/json"
	"github.com/gin-gonic/gin"
)

type fakeToken struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 懒得搞真的token了， 搞个假的算了
func FakeTokenMiddleware(adminConf *AdminConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		var fake fakeToken
		if err := json.ParseJson(token, &fake); err != nil {
			c.JSON(403, gin.H{"success": false, "info": "invalid Token"})
			c.Abort()
			return
		}
		if pass := adminConf.Configs[fake.Username]; pass == "" || fake.Password != pass {
			c.JSON(403, nil)
			c.Abort()
		}
	}
}

func GenFakeToken(username string, adminConf *AdminConfig) string {
	return json.StringifyJson(&fakeToken{
		Username: username,
		Password: adminConf.Configs[username],
	})
}