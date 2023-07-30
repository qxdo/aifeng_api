package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"xiaoniu/common"
	"xiaoniu/module"
	"xiaoniu/utils"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func Route(engine *gin.Engine) {
	engine.POST("/addTask", AddTask)
	engine.POST("/showTask", ShowTask)
	engine.POST("/getTaskById", GetTask)
	engine.GET("/getDetail", ShowDetail)
	engine.POST("/changeStatus", ChangeStatus)
	engine.POST("/deleteTasks", DeleteTasks)
	engine.POST("/setFirst", setFirst)
}

func AddTask(content *gin.Context) {
	task := &common.Task{}
	if err := content.BindJSON(task); err != nil {
		utils.SendResponse(content, 100, "参数错误！", nil)
		println(err.Error())
		return
	}
	fmt.Println("addtask---->", task)
	userId, err := utils.GetId(content)
	if err != nil {
		content.Redirect(302, "/")
		return
	}
	task.Uid = userId
	if task.Type == "" {
		task.Type = "白单"
	}
	title, num, err := common.GetDetailAddTask(task.Url)
	if err != nil {
		utils.SendResponse(content, 500, "获取文章初始量失败！", nil)
		return
	}
	task.BeforeCount = num
	task.Title = title

	// 查询价格
	price, err := module.GetPriceFormDB(task.Uid)
	if err != nil {
		utils.SendResponse(content, 500, "获取价格失败! ", nil)
		return
	}
	if price == nil {
		utils.SendResponse(content, 500, "改用户未设置价格! ", nil)
		return
	}
	if task.Type == "白单" {
		task.Price = float32(price.DayPrice)
	} else {
		task.Price = float32(price.NightPrice)
	}
	id, err := module.AddTask(task)
	if err != nil {
		utils.SendResponse(content, 500, "服务器内部错误！", gin.H{
			"id": id,
		})
		return
	}
	if task.Secret == "---*SECRET_XIAONIU*---" {
		utils.SendResponse(content, 0, "添加成功！", gin.H{
			"before_count": num,
			"id":           id,
		})
	} else {
		utils.SendResponse(content, 0, "添加成功！", nil)
	}
}

func ShowTask(content *gin.Context) {
	uid, err := utils.GetId(content)
	if err != nil {
		content.Redirect(302, "/")
		return
	}
	params := make(map[string]interface{}, 0)
	if err = content.BindJSON(&params); err != nil {
		utils.SendResponse(content, 500, "参数错误！", nil)
		return
	}
	params["uid"] = uid
	result, err := module.ShowTask(params)
	if err != nil {
		utils.SendResponse(content, 500, "服务器内部错误！", nil)
		return
	}
	var param = map[string]interface{}{
		"uid": uid,
	}
	total, err := module.GetTotalTaskNum(param)
	if err != nil {
		utils.SendResponse(content, 500, "服务器内部错误！", nil)
		return
	}
	content.JSON(200, gin.H{
		"code":    0,
		"total":   total,
		"message": "获取成功！",
		"data":    result,
	})
}

func GetTask(content *gin.Context) {
	task := &common.Task{}
	if err := content.BindJSON(task); err != nil {
		utils.SendResponse(content, 100, "参数错误！", nil)
		return
	}
	if task.Id == 0 {
		utils.SendResponse(content, 100, "参数错误！", nil)
		return
	}
	tasks, err := module.GetTaskById(task.Id)
	if err != nil {
		log.Println(err)
		utils.SendResponse(content, 100, "服务器错误！", nil)
		return
	}
	utils.SendResponse(content, 100, "查询成功！", tasks)
}

func ShowDetail(content *gin.Context) {
	uid, err := utils.GetId(content)
	if err != nil {
		content.Redirect(302, "/")
		return
	}
	date := content.Query("date")
	numInfo, err := module.GetNumInfo(uid, date)
	if err != nil {
		utils.SendResponse(content, 500, err.Error(), nil)
		return
	}
	utils.SendResponse(content, 0, "获取成功！", numInfo)
}

func ChangeStatus(content *gin.Context) {
	param := make(map[string]interface{})
	if err := content.BindJSON(&param); err != nil {
		utils.SendResponse(content, 100, "参数错误！", nil)
		return
	}
	paramStr, err := json.Marshal(param)
	if err != nil {
		utils.SendResponse(content, 100, "参数错误！", nil)
		return
	}
	idsRes := gjson.Get(string(paramStr), "ids").Array()
	status := int(gjson.Get(string(paramStr), "status").Int())
	ids := make([]int, 0, len(idsRes))
	for _, id := range idsRes {
		ids = append(ids, int(id.Int()))
	}
	_, err = module.ChangeStatus(ids, status)
	if err != nil {
		fmt.Println(err)
		utils.SendResponse(content, 500, "修改失败！", nil)
		return
	}
	utils.SendResponse(content, 0, "修改成功！", nil)
}

func DeleteTasks(context *gin.Context) {
	param := make(map[string]interface{})
	if err := context.BindJSON(&param); err != nil {
		utils.SendResponse(context, 100, "参数错误！", nil)
		return
	}
	paramStr, err := json.Marshal(param)
	if err != nil {
		utils.SendResponse(context, 100, "参数错误！", nil)
		return
	}
	idsRes := gjson.Get(string(paramStr), "ids").Array()
	ids := make([]int, 0, len(idsRes))
	for _, id := range idsRes {
		ids = append(ids, int(id.Int()))
	}
	_, err = module.DeleteTasks(ids)
	if err != nil {
		fmt.Println(err)
		utils.SendResponse(context, 500, "删除失败！", nil)
		return
	}
	utils.SendResponse(context, 0, "删除成功！", nil)
}

func setFirst(context *gin.Context) {
	param := make(map[string]interface{})
	if err := context.BindJSON(&param); err != nil {
		utils.SendResponse(context, 100, "查询失败，参数错误！", nil)
		return
	}
	var id int
	if temp, ok := param["id"].(float64); ok {
		id = int(temp)
	}
	err := module.SetPriority(id)
	if err != nil {
		utils.SendResponse(context, 500, "设置失败，服务器内部错误！", nil)
		return
	}
	utils.SendResponse(context, 0, "设置成功！", nil)
}
