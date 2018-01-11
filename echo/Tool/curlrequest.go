package Tool

import (
	"io/ioutil"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
)

/**
 *GET 请求
 */
func Get(url string) string {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	//request.Header.Set("Connection", "Keep-Alive")
	//request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	//request.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	//request.Header.Add("Connection", "keep-alive")
	//request.Header.Add("Cookie", "设置cookie")
	//request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	//接收服务端返回给客户端的信息
	response, _ := client.Do(request)
	defer response.Body.Close()
	//cookies := response.Cookies() //遍历cookies
	//for _, cookie := range cookies {
	//	fmt.Println("cookie:", cookie)
	//}
	if response.StatusCode == 200 {
		str, _ := ioutil.ReadAll(response.Body)
		if str == nil {
			return ""
		}
		return string(str)
	} else {
		return ""
	}
}

/**
 * POST 请求
 */
func Post(postUrl string, param map[string]string) string {

	bytesData, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	//fmt.Println(string(bytesData[:]))
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", postUrl, reader)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	//byte数组直接转成string，优化内存
	return string(respBytes)
}
