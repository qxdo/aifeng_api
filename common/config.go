package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
	"xiaoniu/utils"
)

var Conf = &Config{}
var path *string
var LOC, _ = time.LoadLocation("Asia/Shanghai")

const (
	TableTask  = "kpm_read"
	TableProxy = "proxy_sleep"
	TableUser  = "kpm_read_user"
	TablePrice = "kpm_read_price"
)

const SleepCount = 30

var Status = map[int]string{
	0: "全部",
	1: "等待中",
	2: "无效",
	3: "进行中",
	4: "完成",
	5: "补单中",
	6: "停用",
}

const (
	StatusAll = iota
	StatusEnable
	StatusFailed
	StatusRunning
	StatusComplete
	StatusSupplement
	StatusDisabled
)

const (
	NoRowsError = "sql: no rows in result set"
)

var GetNumError = errors.New("cannot get read num")

func CycleLoad() error {
	err := Load()
	if err != nil {
		return err
	}
	conf := Conf.Database
	datasource := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4", conf.Username, conf.Password, conf.Host, conf.Port, conf.Datasource)
	err = utils.InitSqlx(conf.Driver, datasource)
	if err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			err := Load()
			if err != nil {
				fmt.Println("config load error: ", err.Error())
			}
		}
	}()
	return nil
}

func Load() error {
	content, err := ioutil.ReadFile(*path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, Conf)
	if err != nil {
		return err
	}
	return nil
}
