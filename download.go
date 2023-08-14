package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gookit/color"
)

func Download(url string, filepath string) error {
	// HTTP GET请求 下载配置文件
	resp, err := http.Get(url)
	if err != nil {
		color.BgRed.Println("请求失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 将HTTP响应体内容读取到字节切片中
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.BgRed.Println("读取响应体失败:", err)
		return err
	}

	// 将字节切片写入文件
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		color.BgRed.Println("写入文件失败:", err)
		return err
	}

	return nil
}
