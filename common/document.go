package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"xiaoniu/utils"

	"github.com/tidwall/gjson"
)

func GetNum(url string) (int, error) {
	numUrl := fmt.Sprintf("https://api.whosecard.com/api/msg/ext?key=944f5b916ee48bf76e83e6c532ed207201b94fba3a321535a22b8e59&url=%v", url)
	resp, err := utils.Get(numUrl)
	if err != nil {
		println("get num error: ", err.Error())
		return 0, err
	}
	return int(gjson.Get(resp, "clicksCount").Int()), nil
}

// 新修改的方法
func GetDetail(readUrl string) (string, int, error) {
	proxy := Proxy
	if len(proxy) > 100 {
		proxy = proxy[:100]
	}
	for _, v := range proxy {
		url := fmt.Sprintf("http://%s/api/OfficialAccounts/GetAppMsgExt", v.Proxy)
		param := map[string]interface{}{
			"Url":  readUrl,
			"Wxid": v.Guid,
		}
		resp, err := utils.Post(url, param)
		if err != nil {
			fmt.Println("请求出错")
			continue
		}
		res := make(map[string]interface{})
		err = json.Unmarshal([]byte(resp), &res)
		data := fmt.Sprintln(res["Data"])
		code := fmt.Sprintln(res["Code"])
		codes := strings.TrimSpace(code)
		if codes == "0" {
			dataMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(data), &dataMap)
			if _, ok := dataMap["appmsgstat"]; ok {
				appmsgstat := dataMap["appmsgstat"].(map[string]interface{})
				readNum := fmt.Sprintln(appmsgstat["read_num"])
				readNumInt, err := strconv.Atoi(strings.ReplaceAll(readNum, "\n", ""))
				if err != nil {
					fmt.Println("转化出错")
					fmt.Println("错误信息：", err)
				}
				return "", readNumInt, nil
			} else {
				continue
			}

		} else {
			continue
		}

	}
	return "", 0, errors.New("cannot get fullURL")
}

// 已废弃原来的方法名是GetDetail
func GetDetails(readUrl string) (string, int, error) {
	proxy := Proxy
	if len(proxy) > 100 {
		proxy = proxy[:]
	}
	fullUrl, value, err := getFullUrl(readUrl, proxy)

	if err != nil {
		fmt.Println("get fullurl error: ", err)
		return "", 0, err
	}
	title, num, err := dispatchDetail(fullUrl, value)
	if err != nil {
		fmt.Println("get dispatchDetail error: ", err)
		return "", 0, err
	}
	return title, num, nil
}

func GetDetailAddTask(readUrl string) (string, int, error) {
	return GetDetail(readUrl)
	//proxy := Proxy
	//if len(proxy) > 100 {
	//	proxy = proxy[:100]
	//}
	//fullUrl, value, err := getFullUrl(readUrl, proxy)
	//if err != nil {
	//	fmt.Println("get 添加任务err error: ", err)
	//	return "", 0, err
	//}
	//title, num, err := dispatchDetail(fullUrl, value)
	//
	//if err != nil {
	//	return "", 0, err
	//}
	//return title, num, nil
}

func getFullUrl(readUrl string, proxy []*ProxyEntry) (string, string, error) {
	for _, p := range proxy {
		url := fmt.Sprintf("http://%s/api/Common/WXGetA8Key", p.Proxy)
		// url := fmt.Sprintf("http://%s/api/Auth/WXGetMpA8Key", p.Proxy)
		var param = map[string]interface{}{
			"Guid": p.Guid,
			"Url":  readUrl,
		}
		content, err := utils.Post(url, param)
		if err == nil {
			fmt.Println("content=>", content)
			fullUrl := gjson.Get(content, "data").Get("FullURL").String()
			value := gjson.Get(content, "data.HttpHeader.0.Value").String()
			fmt.Println("fullUrl=>", fullUrl)
			if fullUrl == "" {
				fmt.Println("failed to get fullURL retrying")
				continue
			}
			return fullUrl, value, nil
		} else {
			fmt.Println("get detail error: ", err.Error())
		}
	}
	return "", "", errors.New("cannot get fullURL")
}

func dispatchDetail(fullUrl string, value string) (string, int, error) {
	url := "http://49.235.113.82:8012/api/Biz/WXReadBizArticle"
	for i := 0; i < 10; i++ {
		param := map[string]interface{}{
			"Url":   fullUrl,
			"Key":   "exportkey",
			"Value": value,
			"Guid":  "",
		}
		detail, err := utils.Post(url, param)
		if err == nil {
			res := make(map[string]interface{})
			err = json.Unmarshal([]byte(detail), &res)
			var title string
			var num int64
			if code, ok := res["code"].(float64); !ok || code != 0 {
				continue
			}
			if data, ok := res["data"].(string); ok {
				fmt.Println("data-------=>", data)
				title = gjson.Get(data, "Title").String()
				result := gjson.Get(data, "Result").String()
				stat := gjson.Get(result, "appmsgstat").String()
				dataMap := make(map[string]interface{})
				err := json.Unmarshal([]byte(stat), &dataMap)
				if err != nil {
					return "", 0, GetNumError
				}
				if n, ok := dataMap["read_num"].(float64); ok {
					return title, int(n), nil
				}
				fmt.Println("get detail error, json: ", detail)
			}
			if code := gjson.Get(detail, "code").Int(); code == 0 {
				return title, int(num), nil
			}
			return "", 0, errors.New("cannot get detail retrying")
		} else {
			fmt.Println("get detail error: ", err.Error())
		}
	}
	return "", 0, errors.New("cannot get detail")
}

func Read(ch chan *Task, wg *sync.WaitGroup) int {
	defer wg.Done()
	failed := 0
	for {
		task, ok := <-ch
		if !ok {
			return failed
		}
		param := make(map[string]interface{})
		param["Url"] = task.ReadUrl
		param["Wxid"] = task.Guid
		resp, err := utils.Post(task.Url, param)
		if err != nil {
			fmt.Println("get fullURL from proxy error: ", err.Error())
			failed++
			continue
		}
		res := make(map[string]interface{})
		err = json.Unmarshal([]byte(resp), &res)
		code := fmt.Sprintln(res["Code"])
		codes := strings.TrimSpace(code)
		if codes != "0" {
			fmt.Println(fmt.Sprintln(res["Message"]))
			failed++
			continue
		}
		//fullUrl := gjson.Get(resp, "data").Get("FullURL").String()
		//time.Sleep(1)
		//_, err = utils.Get(fullUrl)
		//if err != nil {
		//	fmt.Println("get fullURL error: ", err.Error())
		//	failed++
		//	continue
		//}
	}
}
