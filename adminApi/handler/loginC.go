package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/middleware/jwt"
	"joe-micro/adminApi/model"
	"joe-micro/lib/toolfunc"
)

type login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var res login
	if err := c.ShouldBind(&res); err != nil {
		c.JSON(200, gin.H{
			"code": 3,
			"msg":  fmt.Sprintf("参数错误 %v", err),
		})
		c.Abort()
	}

	//获取用户记录
	userM, has ,err  := model.GetByUsername(res.Username)
	if has {
		c.JSON(200,gin.H{
			"code": 2,
			"msg":  "账号不存在",
		})
		c.Abort()
	}
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("服务器错误 %v", err),
		})
		c.Abort()
	}

    //校验密码
    if   toolfunc.EncUserPwd(res.Password,userM.Salt) != userM.Password {
		c.JSON(200,gin.H{
			"code": 2,
			"msg":  "密码错误",
		})
		c.Abort()
	}

	//颁发token
	token,err:=jwt.CreateToken(userM.ID)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("服务器错误 %v", err),
		})
		c.Abort()
	}


	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": gin.H{
			"token":token,
		},
	})



}
