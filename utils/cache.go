package utils

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
	"typora_uploader_go/logging"

	"golang.org/x/sys/windows/registry"
)

// 下载文件
func CacheFile(imgUrl string, path string) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.QUERY_VALUE)
	if err != nil {
		logging.Logger.Fatal(err)
	}
	val, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		logging.Logger.Fatal(err)
	}
	uri, err := url.Parse("http://" + val)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	client := http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}

	resp, err := client.Get(imgUrl)
	if err != nil {
		logging.Logger.Print(err)
	}
	if resp != nil {
		defer resp.Body.Close()
	} else {
		logging.Logger.Panicf("cannot get '%v'", imgUrl)
	}

	// 创建一个文件用于保存
	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
