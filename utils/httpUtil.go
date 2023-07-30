package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const CodeSuccess = 200

func Post(url string, params map[string]interface{}) (string, error) {
	paramsJson, err := json.Marshal(params)
	if err != nil {
		return "", nil
	}
	client := http.Client{Timeout: 60 * time.Second}
	res, err := client.Post(url,
		"application/json;charset=utf-8", bytes.NewBuffer(paramsJson))
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("http close error: ", err)
		}
	}(res.Body)

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Get(url string) (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("http close error: ", err)
		}
	}(response.Body)
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SendResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(CodeSuccess, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func AddCookie(context *gin.Context, name, value string) {
	context.SetCookie(name, value, 36000, "/", "", false, false)
}

func GetId(context *gin.Context) (int, error) {
	token, err := context.Cookie("token")
	if err != nil || token == "" {
		return 0, err
	}
	id, err := jwt.DecodeSegment(token)
	if err != nil {
		return 0, err
	}
	uid, err := strconv.Atoi(string(id))
	if err != nil {
		return 0, err
	}
	return uid, nil
}

// 判断文章此内容因违规无法查看
func DetermineWhetherViolations(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("请求错误")
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("状态码非200")
		return false
	}
	fmt.Println(res.Body)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("解析html错误")
		return false
	}
	///html/body/div/div[2]/h2
	menu := doc.Find(".tips").Text()
	if menu == "由用户投诉并经平台审核，涉嫌提供未经安全验证的开发者程序下载服务行为，查看对应规则" {
		return true
	}
	return false
}

// 判断文章是否删除
func WhetherArticleDelet(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("请求错误")
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("状态码非200")
		return false
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("解析html错误")
		return false
	}
	menu := doc.Find("div.weui-msg__title").Text()
	if menu == "该内容已被发布者删除" {
		fmt.Println("该内容已被发布者删除")
		return true
	} else {
		return false
	}
}
