package module

import (
	"errors"
	"sync"
	"time"
	"xiaoniu/common"
)

func Flush(url string, num, startNum int, proxy []*common.ProxyEntry, lock *sync.Mutex) (int, error) {
	if len(proxy) == 0 {
		return 0, errors.New("no proxy to flush")
	}
	oriNum := num
	endNum := 0
	for {
		taskChan := make(chan *common.Task, common.Conf.MaxThread)
		wg := &sync.WaitGroup{}
		thread := num
		if thread > common.Conf.MaxThread {
			thread = common.Conf.MaxThread
		}
		lock.Lock()
		for i := 0; i < num; i++ {
			if i < len(proxy) {
				proxy[i].Count++
			}
		}
		lock.Unlock()
		for i := 0; i < thread; i++ {
			wg.Add(1)
			go func() {
				_ = common.Read(taskChan, wg)
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
		time.Sleep(10)
		_, readNum, err := common.GetDetail(url)
		if err != nil {
			return endNum, errors.New("failed to get num")
		}
		endNum = readNum
		num = oriNum - (readNum - startNum)
		if num <= 0 {
			break
		}
		if num >= len(proxy) {
			return endNum, errors.New("insufficient number of agents")
		}
		proxy = proxy[num:]
	}
	return endNum, nil
}
