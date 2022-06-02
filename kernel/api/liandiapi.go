package api

import (
	"github.com/88250/gulu"
	"github.com/gin-gonic/gin"
	"github.com/siyuan-note/siyuan/kernel/conf"
	"github.com/siyuan-note/siyuan/kernel/model"
	"net/http"
)

func getUser(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	model.Conf.UserData = ""
	model.Conf.User = nil
	model.Conf.Save()

	user := conf.User{}
	user.UserId = "1642145846461"
	user.UserTrafficUpload = 121679707
	user.UserCreateTime = "20220114 15:37:26"
	user.UserTrafficDownload = 83477120
	user.UserTrafficTime = 1654121527541
	user.UserTokenExpireTime = "1656726730"
	user.UserSiYuanProExpireTime = -1
	user.UserSiYuanRepoSize = 8000000000
	user.UserAvatarURL = "https://assets.b3logfile.com/avatar/1642145846461.png?imageView2/1/w/256/h/256/interlace/0/q/100"
	user.UserName = "zxcvbnmzsedr"

	user.UserTitles = []*conf.UserTitle{{
		Name: "zxcvbnmzsedr",
		Desc: "desc",
	}}

	data, _ := gulu.JSON.MarshalJSON(user)

	ret.Data = string(data)
}
func serverLogin(c *gin.Context) {
	ret := &gulu.Result{
		Code: 0,
		Msg:  "",
		Data: map[string]interface{}{
			"userName":    "zxcvbnmzsedr",
			"token":       "123",
			"needCaptcha": false,
		},
	}

	defer c.JSON(http.StatusOK, ret)
}
