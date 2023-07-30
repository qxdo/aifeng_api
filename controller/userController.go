package controller

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaoniu/common"
	"xiaoniu/module"
	"xiaoniu/utils"
)

func UserRote(engine *gin.Engine) {
	engine.Handle("POST", "/login", login)
	engine.Handle("GET", "/checkLogin", checkLogin)
}

func login(context *gin.Context) {
	user := common.User{}
	if err := context.BindJSON(&user); err != nil {
		utils.SendResponse(context, 100, "参数不正确！", nil)
		return
	}
	u, err := module.Login(user)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendResponse(context, 200, "用户名或密码错误！", nil)
		} else {
			utils.SendResponse(context, 500, "发生未知错误，登陆失败！", nil)
		}
		return
	}
	token := jwt.EncodeSegment([]byte(strconv.Itoa(u.Id)))
	utils.AddCookie(context, "token", token)
	utils.SendResponse(context, 0, "登陆成功！", nil)
}

func checkLogin(context *gin.Context) {
	_, err := utils.GetId(context)
	if err != nil {
		utils.SendResponse(context, 100, "请登录！", nil)
		return
	}
	utils.SendResponse(context, 0, "令牌通过！", nil)
}
