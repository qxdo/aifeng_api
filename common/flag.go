package common

import "flag"

func CommandLine() {
	path = flag.String("config", "config/config.json", "指定配置文件")
	flag.Parse()
}
