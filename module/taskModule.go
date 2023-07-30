package module

import (
	"fmt"
	"strconv"
	"time"
	"xiaoniu/common"
	"xiaoniu/utils"

	"github.com/jmoiron/sqlx"
)

func AddTask(task *common.Task) (int, error) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	sqlStr := "insert into " + common.TableTask + " (uid,title,url,`type`,price,demand_count,before_count,all_count,status,add_time,end_time,is_first) " +
		"values (?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := utils.Execute(sqlStr, task.Uid, task.Title, task.Url, task.Type, task.Price, task.DemandCount,
		task.BeforeCount, task.AllCount, common.StatusEnable, currentTime, currentTime, 1)
	id, err := (*result).LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func ShowTask(params map[string]interface{}) ([]*common.Task, error) {
	var limit, page int
	if val, ok := params["page"].(float64); ok {
		page = int(val)
		delete(params, "page")
	}
	if val, ok := params["limit"].(float64); ok {
		limit = int(val)
		delete(params, "limit")
	}
	values := make([]interface{}, 0, len(params))
	sqlCmd := fmt.Sprintf("select id,url,uid,title,price,suc_count,demand_count,before_count,status,add_time,`type`,end_time from %s where is_delete=0 and ", common.TableTask)
	if val, ok := params["time"].(string); ok {
		sqlCmd += fmt.Sprintf("DATE_FORMAT(add_time, '%%Y-%%m-%%d') = '%s' and ", val)
		delete(params, "time")
	}
	if val, ok := params["title"].(string); ok {
		sqlCmd += fmt.Sprintf("(title like '%%%s%%' or url='%s') and ", val, val)
		delete(params, "title")
	}
	for key, value := range params {
		switch val := value.(type) {
		case string:
			if val == "全部" {
				continue
			}
		case float64:
			if val == 0.0 {
				continue
			}
		case int:
			if val == 0 {
				continue
			}
		case nil:
			continue
		}
		sqlCmd += "`" + key + "`" + "=? and "
		values = append(values, value)
	}
	sqlCmd = sqlCmd[:len(sqlCmd)-5]
	sqlCmd += " order by add_time desc"
	if limit != 0 && page != 0 {
		sqlCmd += fmt.Sprintf(" limit %d,%d", (page-1)*limit, limit)
	}
	rows, err := utils.QueryNormal(sqlCmd, values...)
	if err != nil {
		return nil, err
	}
	tasks := make([]*common.Task, 0)
	for rows.Next() {
		task := common.Task{}
		s := 0
		err := rows.Scan(&task.Id, &task.Url, &task.Guid, &task.Title, &task.Price, &task.SucCount, &task.DemandCount, &task.BeforeCount, &s, &task.AddTime, &task.Type, &task.EndTime)
		if err != nil {
			fmt.Println("scan error: ", err.Error())
		}
		task.Status = common.Status[s]
		tasks = append(tasks, &task)
	}
	err = rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	return tasks, nil
}

func GetTotalTaskNum(param map[string]interface{}) (int, error) {
	sqlCmd := fmt.Sprintf("select count(*) from %s where ", common.TableTask)
	values := make([]interface{}, 0, len(param))
	for key, value := range param {
		values = append(values, value)
		sqlCmd += key + "=? and "
	}
	sqlCmd = sqlCmd[:len(sqlCmd)-5]
	row := utils.DB.QueryRow(sqlCmd, values...)
	num := 0
	err := row.Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func GetTaskById(id int) ([]*common.Task, error) {
	sqlStr := fmt.Sprintf("select id,url,uid,title,price,suc_count,all_count,status,add_time,`type`,end_time from %s where id=?", common.TableTask)
	rows, err := utils.DB.Queryx(sqlStr, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("db close error: ", err.Error())
		}
	}(rows)
	tasks := make([]*common.Task, 0)
	for rows.Next() {
		task := common.Task{}
		err := rows.Scan(&task.Id, &task.Url, &task.Guid, &task.Title, &task.Price, &task.SucCount, &task.AllCount, &task.Status, &task.AddTime, &task.Type, &task.EndTime)
		if err != nil {
			fmt.Println("scan error: ", err.Error())
		}
		tasks = append(tasks, &task)
	}
	err = rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	return tasks, nil
}

func GetNumInfo(uid int, queryDate string) (*common.CountInfo, error) {
	t := time.Now()
	date := t.Format("2006-01-02")
	yt := t.AddDate(0, 0, -1)
	yestDay := yt.Format("2006-01-02")
	numInfo := common.CountInfo{}
	proxy := common.Proxy
	// 计算代理数量
	numInfo.AllAccount = len(proxy)
	for _, p := range proxy {
		if p.Count >= common.SleepCount {
			numInfo.SleepAccount++
		} else {
			numInfo.EnableAccount++
		}
	}
	// 计算任务数量
	sqlCmd := ""
	var row *sqlx.Row
	if queryDate == "" {
		sqlCmd = fmt.Sprintf("select ifnull(count(*),0) as task_num,"+ // 任务总数
			"(select ifnull(count(*),0) from %s where status=%d and uid=?) as running_num,"+ // 进行中
			"(select ifnull(count(*),0) from %s where status=%d and uid=?) as completed_num,"+ // 已完成
			"(select ifnull(sum(demand_count),0) from %s where is_delete=0 and uid=?) as task_count,"+
			"ifnull(sum(suc_count),0) as completed_count,ifnull(sum(price*suc_count),0) as total_price "+ // 任务总量、完成总量、总价
			"from %s where is_delete=0 and uid=?",
			common.TableTask, common.StatusRunning, common.TableTask, common.StatusComplete, common.TableTask, common.TableTask)
		row = utils.DB.QueryRowx(sqlCmd, uid, uid, uid, uid)
	} else {
		queryDate = fmt.Sprintf("%%%s%%", queryDate)
		sqlCmd = fmt.Sprintf("select ifnull(count(*),0) as task_num,"+ // 任务总数
			"(select ifnull(count(*),0) from %s where status=%d and uid=? and add_time like ?) as running_num,"+ // 进行中
			"(select ifnull(count(*),0) from %s where status=%d and uid=? and add_time like ? and is_delete=0) as completed_num,"+ // 已完成
			"(select ifnull(sum(demand_count),0) from %s where is_delete=0 and uid=? and add_time like ?) as task_count,"+
			"ifnull(sum(suc_count),0) as completed_count,ifnull(sum(price*suc_count),0) as total_price "+ // 任务总量、完成总量、总价
			"from %s where is_delete=0 and uid=? and add_time like ?",
			common.TableTask, common.StatusRunning, common.TableTask, common.StatusComplete, common.TableTask, common.TableTask)
		row = utils.DB.QueryRowx(sqlCmd, uid, queryDate, uid, queryDate, uid, queryDate, uid, queryDate)
	}

	err := row.StructScan(&numInfo)
	if err != nil {
		fmt.Println(err)
		return &numInfo, nil
	}
	sqlCmd = fmt.Sprintf("select (select ifnull(sum(demand_count),0) from %s where `type`='白单' and is_delete=0 and uid=? and date(add_time)=?) as day_count,"+
		"(select ifnull(sum(demand_count),0) from %s where `type`='夜单' and uid=? and is_delete=0 and date(add_time)=?) as night_count,"+
		"(select ifnull(sum(demand_count),0) from %s where `type`='夜单' and uid=? and is_delete=0 and date(add_time)=?) as yest_night_count,"+
		"(select ifnull(sum(demand_count),0) from %s where `type`='白单' and is_delete=0 and uid=? and date(add_time)=?) as yest_day_count",
		common.TableTask, common.TableTask, common.TableTask, common.TableTask)
	row = utils.DB.QueryRowx(sqlCmd, uid, date, uid, date, uid, yestDay, uid, yestDay)
	err = row.StructScan(&numInfo)
	if err != nil {
		return nil, err
	}
	return &numInfo, nil
}

func Recover() error {
	sqlCmd := fmt.Sprintf("update %s set status=1 where status=3", common.TableTask)
	ret, err := utils.Execute(sqlCmd)
	if err != nil {
		return err
	}
	rowNum, err := (*ret).RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("recover rows: ", rowNum)
	sqlCmd = fmt.Sprintf("update %s set status=4 where suc_count>=demand_count", common.TableTask)
	ret, err = utils.Execute(sqlCmd)
	if err != nil {
		return err
	}
	rowNum, err = (*ret).RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("recover status success rows: ", rowNum)
	return nil
}

func ChangeStatus(ids []int, status int) (int, error) {
	if len(ids) < 1 {
		return 0, nil
	}
	idStr := strconv.Itoa(ids[0])
	for i := 1; i < len(ids); i++ {
		idStr += "," + strconv.Itoa(ids[i])
	}
	sqlCmd := fmt.Sprintf("update %s set status=%d where id in (%s) and disabled!=1", common.TableTask, status, idStr)
	result, err := utils.Execute(sqlCmd)
	if err != nil {
		return 0, err
	}
	num, err := (*result).RowsAffected()
	if err != nil {
		return 0, nil
	}
	return int(num), nil
}

func DeleteTasks(ids []int) (int, error) {
	if len(ids) < 1 {
		return 0, nil
	}
	idStr := strconv.Itoa(ids[0])
	for i := 1; i < len(ids); i++ {
		idStr += "," + strconv.Itoa(ids[i])
	}
	sqlCmd := fmt.Sprintf("update %s set is_delete=1 where id in (%s)", common.TableTask, idStr)
	result, err := utils.Execute(sqlCmd)
	if err != nil {
		return 0, err
	}
	num, err := (*result).RowsAffected()
	if err != nil {
		return 0, nil
	}
	return int(num), nil
}

func SetPriority(id int) error {
	sqlCmd := fmt.Sprintf("update %s set priority=1 where id=?", common.TableTask)
	_, err := utils.Execute(sqlCmd, id)
	if err != nil {
		return err
	}
	return nil
}
