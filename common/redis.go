package common

import "xiaoniu/utils"

const (
	_Host     = "127.0.0.1"
	_Port     = 6379
	_DB       = 13
	_Password = "922NVNsdAhGEW5pP"
)

func init() {
	err := utils.InitClient(_Host, _Port, _Password, _DB)
	if err != nil {
		panic("redis init error: " + err.Error())
	}
}
