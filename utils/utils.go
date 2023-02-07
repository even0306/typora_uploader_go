package utils

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"typora_uploader_go/logs"

	"github.com/google/uuid"
)

var logging = logs.LogFile()

// 判断数据类型
func FileType(file *string) (filetype string) {

	if *file != "" {
		req := strings.Split(*file, ":")[0]
		if req == "data" {
			filetype = "base64"
		} else if req == "http" || req == "https" {
			filetype = "url"
		} else {
			filetype = "local"
		}
	} else {
		logging.Printf("数据不能为空")
	}
	return
}

// 读取文件为二进制格式
func ReadFile(path *string) (b *[]byte, e error) {
	file, err := os.Open(*path)
	if err != nil {
		logging.Printf("打开文件失败, error: %v", err)
		return
	}
	defer file.Close()
	chunks := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			logging.Printf("读取文件失败，error: %v", err)
		}
		if n == 0 {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return &chunks, err
}

// 下载文件
func DownloadFile(imgUrl *string, path *string) {
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(*imgUrl)
	if err != nil {
		logging.Print(err)
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(*path)
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

// 创建UUID作为文件名
func CreateUUID() (key string) {
	uuid := uuid.New()
	key = uuid.String()
	return
}
