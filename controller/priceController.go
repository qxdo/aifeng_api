package controller

import (
	"github.com/gin-gonic/gin"
	"time"
	"xiaoniu/common"
	"xiaoniu/module"
	"xiaoniu/utils"
)

func PriceRoute(engine *gin.Engine) {
	engine.Handle("POST", "/setPrice", setPrice)
	engine.Handle("GET", "/getPrice", getPrice)
}

func getPrice(context *gin.Context) {
	uid, err := utils.GetId(context)
	if err != nil {
		utils.SendResponse(context, 100, "请登录!", nil)
		return
	}
	price, err := module.GetPriceFormDB(uid)
	if err != nil {
		utils.SendResponse(context, 500, "get price failed", nil)
		return
	}
	if price == nil {
		utils.SendResponse(context, 100, "not have price data", nil)
		return
	}
	utils.SendResponse(context, 0, "success", gin.H{
		"day_price":   int(price.DayPrice * 1000),
		"night_price": int(price.NightPrice * 1000),
		"time":        price.SqlTime,
	})
	return
}

func setPrice(context *gin.Context) {
	uid, err := utils.GetId(context)
	if err != nil {
		utils.SendResponse(context, 100, "等登录！", nil)
		return
	}
	params := make(map[string]interface{})
	if err := context.BindJSON(&params); err != nil {

	}
	var validDayPrice bool
	var validNightPrice bool
	price := &common.Price{}
	price.Uid = uid
	if dp, ok := params["dayPrice"].(float64); ok {
		price.DayPrice = dp / 1000
		validDayPrice = true
	}
	if np, ok := params["nightPrice"].(float64); ok {
		price.NightPrice = np / 1000
		validNightPrice = true
	}
	if !(validDayPrice && validNightPrice) {
		utils.SendResponse(context, 100, "参数错误！", nil)
		return
	}
	st := time.Now().Format("2006-01-02 15:04:05")
	price.SqlTime = st

	if err := module.AddPriceFormDB(*price); err != nil {
		utils.SendResponse(context, 500, "Add price failed", nil)
		return
	}
	utils.SendResponse(context, 0, "success", nil)
	return
}
