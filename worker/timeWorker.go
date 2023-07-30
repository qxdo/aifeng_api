package worker

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"xiaoniu/common"
	"xiaoniu/module"
	"xiaoniu/utils"
)

type Worker struct {
	taskChan  chan *common.Task
	lock      *sync.Mutex
	wg        *sync.WaitGroup
	timeStamp string
}

func CreateWorker() *Worker {
	taskChan := make(chan *common.Task)
	wg := &sync.WaitGroup{}
	currentDate := time.Now().Format("2006-01-02")
	return &Worker{
		taskChan:  taskChan,
		wg:        wg,
		lock:      &sync.Mutex{},
		timeStamp: currentDate,
	}
}

func (worker *Worker) Start() {
	for {
		worker.createTask()
		currentDate := time.Now().Format("2006-01-02")
		if worker.timeStamp != currentDate {
			err := module.ResetProxyCount()
			if err != nil {
				fmt.Println("reset proxy error: ", err.Error())
			}
			worker.timeStamp = currentDate
		}
	}
}

func (worker *Worker) createTask() {
	defer time.Sleep(2 * time.Second)
	runTime := time.Now().Format("2006-01-02")
	runTime += fmt.Sprintf(" 23:59:59")
	sleepTime := time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
	sqlCmd := fmt.Sprintf("select * from %s "+
		"where is_delete!=1 and "+
		"(status=1 and (type='白单' or type='夜单' and add_time<?) "+
		"or status=5 and end_time<?) "+
		"order by priority desc,status,add_time limit ?",
		common.TableTask)
	rows, err := utils.DB.Queryx(sqlCmd, runTime, sleepTime, common.Conf.TaskNum)
	if err != nil {
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("rows close error: ", err.Error())
		}
	}()
	// 唤醒代理
	err = module.AwakeProxy()
	if err != nil {
		fmt.Println("proxy awake error: ", err.Error())
		return
	}
	proxy, err := module.GetProxy()
	if err != nil {
		fmt.Println("get awake proxy error: ", err)
		return
	}
	if len(proxy) == 0 {
		time.Sleep(2 * time.Minute)
		return
	}
	err = module.GetProxyCount(proxy)
	if err != nil {
		fmt.Println("get awake proxy error: ", err)
		return
	}
	common.Proxy = proxy
	proxy = make([]*common.ProxyEntry, 0, len(proxy))
	for _, p := range common.Proxy {
		if p.Count < common.SleepCount {
			proxy = append(proxy, p)
		}
	}
	ids := make([]string, 0, common.Conf.TaskNum)
	for rows.Next() {
		task := common.Task{}
		err := rows.StructScan(&task)
		if err != nil || task.DemandCount == 0 {
			continue
		}
		ids = append(ids, strconv.Itoa(task.Id))
		worker.wg.Add(1)
		go worker.Do(&task, proxy)

		go func() {
			if utils.WhetherArticleDelet(task.Url) || utils.DetermineWhetherViolations(task.Url) {
				updateSql := fmt.Sprintf("UPDATE %s SET is_delete=1,status=6  WHERE id = ?", common.TableTask)
				_, err := utils.DB.Exec(updateSql, task.Id)
				if err != nil {
					fmt.Println("内容被删除修改数据失败")
				}
			}
		}()
		//if utils.WhetherArticleDelet(task.Url) || utils.DetermineWhetherViolations(task.Url) {
		//	updateSql := fmt.Sprintf("UPDATE %s SET is_delete=1,status=6  WHERE id = ?", common.TableTask)
		//	result, err := utils.DB.Exec(updateSql, task.Id)
		//	if err != nil {
		//		fmt.Println("内容被删除修改数据失败")
		//	}
		//	fmt.Println(result)
		//	continue
		//} else {
		//	ids = append(ids, strconv.Itoa(task.Id))
		//	worker.wg.Add(1)
		//	go worker.Do(&task, proxy)
		//}

	}
	if len(ids) == 0 {
		time.Sleep(10 * time.Second)
		return
	}
	idStr := strings.Join(ids, ",")
	fmt.Println(idStr)
	sqlCmd = fmt.Sprintf("update %s set status=3 where id in (%s)", common.TableTask, idStr)
	_, err = utils.DB.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err)
	}
	err = rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	worker.wg.Wait()
	runTime = time.Now().Format("2006-01-02 15:04:05")
	sqlCmd = fmt.Sprintf("insert into %s (guid,proxy,count,time) values ", common.TableProxy)
	sign := true
	values := make([]interface{}, 0, 4*len(proxy))
	idArr := make([]interface{}, 0, len(proxy))
	idStr = ""
	for _, p := range proxy {
		if p.Count > 0 {
			if sign {
				sign = false
				idStr = fmt.Sprintf("?")
				sqlCmd += fmt.Sprintf("(?,?,?,?)")
			} else {
				idStr += fmt.Sprintf(",?")
				sqlCmd += fmt.Sprintf(",(?,?,?,?)")
			}
			values = append(values, p.Guid, p.Proxy, p.Count, runTime)
			idArr = append(idArr, p.Guid)
		}
	}
	if len(idArr) == 0 {
		return
	}
	delSql := fmt.Sprintf("delete from %s where guid in (%s)", common.TableProxy, idStr)
	_, err = utils.DB.Exec(delSql, idArr...)
	if err != nil {
		fmt.Println("insert proxy error: ", err)
		return
	}
	_, err = utils.DB.Exec(sqlCmd, values...)
	if err != nil {
		fmt.Println("insert proxy error: ", err)
		return
	}
}

func (worker *Worker) Do(task *common.Task, proxy []*common.ProxyEntry) {
	defer worker.wg.Done()
	// 获取初始量
	_, startNum, err := common.GetDetail(task.Url)
	if err != nil {
		fmt.Println("get num error: ", err)
		sqlCmd := fmt.Sprintf("update %s set status=? where id=?", common.TableTask)
		_, err := utils.Execute(sqlCmd, common.StatusSupplement, task.Id)
		if err != nil {
			fmt.Println("sql execute error: ", err)
		}
		return
	}
	// 判断成功量
	task.SucCount = startNum - task.BeforeCount
	if task.SucCount > task.DemandCount {
		task.SucCount = task.DemandCount
		sqlCmd := fmt.Sprintf("update %s set status=? where id=?", common.TableTask)
		_, err := utils.Execute(sqlCmd, common.StatusComplete, task.Id)
		if err != nil {
			fmt.Println("update task error: ", err)
			return
		}
	}
	// 开始处理任务
	num := task.DemandCount - task.SucCount
	readNum, err := module.Flush(task.Url, num, startNum, proxy, worker.lock)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	if err != nil {
		fmt.Println(err)
	}
	sqlCmd := ""
	if readNum != 0 {
		task.SucCount = readNum - task.BeforeCount
	}
	status := common.StatusSupplement
	if task.SucCount >= task.DemandCount {
		task.SucCount = task.DemandCount
		status = common.StatusComplete
	}
	sqlCmd = fmt.Sprintf("update %s set status=?,suc_count=?,is_first=0,end_time=? where id=?", common.TableTask)
	_, err = utils.Execute(sqlCmd, status, task.SucCount, currentTime, task.Id)
	if err != nil {
		fmt.Println("FATAL: update task error: ", err)
	}
}
