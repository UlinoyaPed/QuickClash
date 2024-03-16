package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

func DownloadWithBar(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create our progress bar
	bar := progressbar.Default(resp.ContentLength)

	// Create our reader and stat the size of our file
	reader := io.TeeReader(resp.Body, bar)

	// Write the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	bar.Add(int(resp.ContentLength))
	bar.Finish()

	return nil
}

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

func AddHTTPSPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return url
}
