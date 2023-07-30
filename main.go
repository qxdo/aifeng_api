package main

import (
	"fmt"
	"sync"
	"time"
	"xiaoniu/common"
	"xiaoniu/controller"
	"xiaoniu/module"
	"xiaoniu/worker"

	"github.com/gin-gonic/gin"
)

func flush(c *gin.Context) {
	p := &common.Params{}
	if err := c.BindJSON(p); err != nil {
		c.JSON(200, gin.H{
			"code":    100,
			"message": "参数错误！",
		})
		println(err)
		return
	}
	url := p.Url
	num := p.Num
	startNum := -1
	var err error
	proxy := common.Proxy
	_, startNum, err = common.GetDetail(url)
	if err != nil {
		fmt.Println("failed to get num", err.Error())
	}
	oriNum := num
	endNum := 0
	lock := &sync.RWMutex{}
	for {
		failedNum := 0
		taskChan := make(chan *common.Task, common.Conf.MaxThread)
		wg := &sync.WaitGroup{}
		thread := num
		if thread > common.Conf.MaxThread {
			thread = common.Conf.MaxThread
		}
		for i := 0; i < thread; i++ {
			wg.Add(1)
			go func() {
				failed := common.Read(taskChan, wg)
				lock.Lock()
				failedNum += failed
				lock.Unlock()
			}()
		}
		for i := 0; i < num; i++ {
			if i < len(proxy) {
				t := &common.Task{
					Url: "http://" + proxy[i].Proxy + "/api/OfficialAccounts/GetAppMsgExt",
					//Url: "http://" + proxy[i].Proxy + "/api/Common/WXGetA8Key",
					// Url:     "http://" + proxy[i].Proxy + "/api/Auth/WXGetMpA8Key",
					ReadUrl: url,
					Guid:    proxy[i].Guid,
				}
				taskChan <- t
			}
		}
		close(taskChan)
		wg.Wait()
		_, readNum, err := common.GetDetail(url)
		endNum = readNum
		if err != nil {
			c.JSON(200, gin.H{
				"code":    500,
				"message": "failed to get num",
			})
			return
		}
		num = oriNum - (readNum - startNum)
		if num <= 0 {
			break
		}
		time.Sleep(5)
		if num >= len(proxy) {
			c.JSON(200, gin.H{
				"code":    500,
				"message": "insufficient number of agents",
			})
			return
		}
		proxy = proxy[num:]
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "成功！",
		"start":   startNum,
		"end":     endNum,
	})
}

func route(r *gin.Engine) {
	r.POST("/flush", flush)
}

func main() {
	common.CommandLine()
	err := common.CycleLoad()
	if err != nil {
		panic("config load error: " + err.Error())
	}
	module.InitProxy()
	err = module.Recover()
	if err != nil {
		fmt.Println("task recover failed with error: ", err)
	}
	taskWorker := worker.CreateWorker()
	go taskWorker.Start()
	r := gin.Default()
	route(r)
	controller.Route(r)
	controller.UserRote(r)
	controller.PriceRoute(r)
	err = r.Run(":8080")
	if err != nil {
		fmt.Println("server start failed, err: ", err.Error())
	}
}
