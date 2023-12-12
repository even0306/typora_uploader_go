package platform

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
	"typora_uploader_go/logging"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func UploadSelecter(pt MyLocal, b *[]byte) string {
	switch pt.PicBed {
	case "nextcloud":
		viewURL, err := NextcloudUploadFile(pt, b)
		if err != nil {
			logging.Logger.Panicf("未获取到下载地址：%v", err)
		}
		return viewURL
	case "aliyunOss", "minIO":
		viewURL := OssUploadFile(pt, b)
		return viewURL
	default:
		logging.Logger.Println("不支持的平台")
	}
	return ""
}

// TimeoutDialer 连接超时和传输超时
func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

// 上传接口，传url，文件二进制，参数头
func NextcloudUploadFile(header MyLocal, fileByte *[]byte) (string, error) {
	auth := map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(header.AccessKeyId+":"+header.AccessKeySecret))}
	req, err := http.NewRequest("PUT", header.UploadURL, bytes.NewBuffer(*fileByte))
	if err != nil {
		logging.Logger.Printf("http newrequest error %s", err)
		return "", err
	}
	for h, v := range auth {
		req.Header.Set(h, v)
	}

	connectTimeout := 120 * time.Second
	readWriteTimeout := 5184000 * time.Millisecond
	client := http.Client{
		//忽略证书验证
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			Dial: timeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		logging.Logger.Printf("http client error %s", err)
		return "", err
	}
	if resp != nil {
		// 判断请求状态
		if resp.StatusCode == http.StatusOK {
			respData, err := io.ReadAll(resp.Body)
			if err != nil {
				logging.Logger.Print(err)
				return "", err
			}
			logging.Logger.Printf("\n【请求地址】： %s \n【请求参数】： %s \n【请求头】： %s \n【返回】 : %s \n",
				header.UploadURL, "上传文件", auth, string(respData))
			fmt.Println(string(respData))
			return "", nil
		} else if resp.StatusCode != http.StatusOK {
			respData, err := io.ReadAll(resp.Body)
			if err != nil {
				logging.Logger.Print(err)
				return "", err
			}

			logging.Logger.Printf("\n【请求地址】： %s \n【请求参数】： %s \n【请求头】： %s \n【返回】 : %s \n", header.UploadURL, "上传文件", auth, string(respData))
			logging.Logger.Println("上传文件请求成功，上传成功")
			fmtUrl := header.DownloadURL + "\n"
			return fmtUrl, nil
		}
		return "", errors.New("请求失败")
	}
	defer resp.Body.Close()

	return "", errors.New("请求失败")
}

// 阿里云OSS上传
func OssUploadFile(header MyLocal, fileByte *[]byte) string {
	// func AliyunOssUploadFile(endpoint *string, accessKeyId *string, accessKeySecret *string, bucketName *string, fileName string, fileByte *[]byte) string {
	client, err := oss.New(header.UploadURL, header.AccessKeyId, header.AccessKeySecret)
	if err != nil {
		// HandleError(err)
		logging.Logger.Printf("创建阿里云上传客户端失败，error：%v", err)
		os.Exit(-1)
	}

	bucket, err := client.Bucket(header.BucketName)
	if err != nil {
		// HandleError(err)
		logging.Logger.Printf("获取阿里云桶失败，error：%v", err)
		os.Exit(-1)
	}

	err = bucket.PutObject(header.FileName, bytes.NewReader([]byte(*fileByte)))
	if err != nil {
		// HandleError(err)
		logging.Logger.Printf("阿里云上传失败，error：%v", err)
		os.Exit(-1)
	}
	return "https://" + header.DownloadURL
}
