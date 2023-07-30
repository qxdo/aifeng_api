package module

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
	"xiaoniu/common"
	"xiaoniu/utils"

	"github.com/jmoiron/sqlx"
)

func ResetProxyCount() error {
	sqlCmd := fmt.Sprintf("delete from %s", common.TableProxy)
	_, err := utils.Execute(sqlCmd)
	if err != nil {
		return err
	}
	return nil
}

func AwakeProxy() error {
	t := time.Now().Add(-60 * time.Minute)
	sqlCmd := fmt.Sprintf("delete from %s where time<?", common.TableProxy)
	_, err := utils.Execute(sqlCmd, t.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return nil
}

func GetProxyFormDb(status int) ([]*common.ProxyEntry, error) {
	where := fmt.Sprintf(" where count%s%d", "%s", common.SleepCount)
	switch status {
	case common.All:
		where = ""
	case common.Sleep:
		where = fmt.Sprintf(where, ">=")
	case common.Awake:
		where = fmt.Sprintf(where, "<")
	}
	sqlCmd := fmt.Sprintf("select * from %s", common.TableProxy)
	sqlCmd += where
	rows, err := utils.DB.Queryx(sqlCmd)
	if err != nil {
		return nil, nil
	}
	proxy := make([]*common.ProxyEntry, 0)
	for rows.Next() {
		p := &common.ProxyEntry{}
		err := rows.StructScan(p)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		proxy = append(proxy, p)
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("db close error: ", err.Error())
		}
	}(rows)
	return proxy, nil
}

func GetProxyCount(proxy []*common.ProxyEntry) error {
	dbProxy, err := GetProxyFormDb(common.All)
	if err != nil {
		return err
	}
	proxyMap := make(map[string]int)
	for _, p := range dbProxy {
		proxyMap[p.Guid] = p.Count
	}
	for _, p := range proxy {
		if val, ok := proxyMap[p.Guid]; ok {
			p.Count = val
		}
	}
	return nil
}

func InitProxy() {
	proxy, err := GetProxy()
	if err != nil {
		panic(err)
	}
	err = GetProxyCount(proxy)
	if err != nil {
		panic(err)
	}
	common.Proxy = proxy
	fmt.Println("proxy init success!")
}

func GetProxy() ([]*common.ProxyEntry, error) {
	type task struct {
		key   string
		proxy *common.ProxyEntry
	}
	taskChan := make(chan task, 5000)
	wg := &sync.WaitGroup{}
	keys, err := utils.Redis.GetKeys("*")
	if err != nil {
		return nil, nil
	}
	proxy := make([]*common.ProxyEntry, 0, len(keys))
	// 创建消费z者
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				t, ok := <-taskChan
				if !ok {
					break
				}
				//不检查在线代理
				//p, key := t.proxy, t.key
				proxy = append(proxy, t.proxy)
				//if !CheckAlive(p.Proxy, p.Guid) {
				//	if err := utils.Redis.DelItem(key); err != nil {
				//		fmt.Println("remove redis error: ", err.Error())
				//	}
				//} else {
				//	proxy = append(proxy, p)
				//}
			}
		}()
	}
	for _, key := range keys {
		lines, err := utils.Redis.GetString(key)
		if err != nil {
			fmt.Println("出现错误，key=>", key)
			fmt.Println("get proxy from redis error: ", err.Error())
			continue
		}
		//list变string
		temp := strings.Split(lines, "@")
		if len(temp) == 2 {
			p := &common.ProxyEntry{
				Proxy: temp[0],
				Guid:  temp[1],
			}
			t := task{
				proxy: p,
				key:   key,
			}
			taskChan <- t
		} else {
			err := utils.Redis.DelItem(key)
			if err != nil {
				fmt.Println("proxy remove error: ", err.Error())
			}
		}
		//for _, line := range lines {
		//	temp := strings.Split(line, "@")
		//	if len(temp) == 2 {
		//		p := &common.ProxyEntry{
		//			Proxy: temp[0],
		//			Guid:  temp[1],
		//		}
		//		t := task{
		//			proxy: p,
		//			key:   key,
		//		}
		//		taskChan <- t
		//	} else {
		//		err := utils.Redis.RemoveItem(key, line)
		//		if err != nil {
		//			fmt.Println("proxy remove error: ", err.Error())
		//		}
		//	}
		//}
	}
	close(taskChan)
	wg.Wait()
	fmt.Println(fmt.Sprintf("proxy list flushed, time: %s", time.Now().Format("2006-01-02 15:04:05")))
	return proxy, nil
}

func CheckAlive(ip, guid string) bool {
	url := fmt.Sprintf(common.CheckUrlIpad, ip, guid)
	param := map[string]interface{}{
		"wxid": guid,
	}
	content, err := utils.Post(url, param)
	if err != nil {
		return false
	}
	resp := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &resp)
	if err != nil {
		return false
	}
	if retCode, ok := resp["Code"].(float64); ok && retCode == 0 {
		return true
	}
	//if data, ok := resp["Data"].(map[string]interface{}); ok {
	//	fmt.Println(ok)
	//	if retCode, ok := data["Ret"].(float64); ok && retCode == 0 {
	//		return true
	//	}
	//}
	return false
}
