package main

import (
	"fmt"
	"xiaoniu/utils"
)

func NewFunc() {
	ok := utils.DetermineWhetherViolations("https://cloud.tencent.com/developer/article/2064619")
	fmt.Println(ok)
}
func main() {
	NewFunc()
}
