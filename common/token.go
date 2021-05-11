package common

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type fakeToken struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 懒得搞真的token了， 搞个假的算了
func FakeTokenMiddleware(adminConf *AdminConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		slice := strings.Split(token, "_")
		var fake = fakeToken{
			Username: slice[0],
			Password: slice[len(slice)-1],
		}
		if pass := adminConf.Configs[fake.Username]; pass == "" || fake.Password != pass {
			c.JSON(403, nil)
			c.Abort()
		}
	}
}

func GenFakeToken(username string, adminConf *AdminConfig) string {
	//return json.StringifyJson(&fakeToken{
	//	Username: username,
	//	Password: adminConf.Configs[username],
	//})
	return username + "_" + adminConf.Configs[username]
}
