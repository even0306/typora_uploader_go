package upload

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
	"typora_uploader_go/logs"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

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

var logging = logs.LogFile()

// 上传接口，传url，文件二进制，参数头
func NextcloudUploadFile(rURL string, b *[]byte, header *map[string]string) error {

	req, err := http.NewRequest("PUT", rURL, bytes.NewBuffer(*b))
	if err != nil {
		logging.Printf("http newrequest error %s", err)
		return err
	}
	for h, v := range *header {
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
		logging.Printf("http client error %s", err)
		return err
	}
	if resp != nil {
		// 判断请求状态
		if resp.StatusCode == http.StatusOK {
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logging.Print(err)
				return err
			}
			logging.Printf("\n【请求地址】： %s \n【请求参数】： %s \n【请求头】： %s \n【返回】 : %s \n",
				rURL, "上传文件", *header, string(respData))
			fmt.Println(string(respData))
			return nil
		} else if resp.StatusCode != http.StatusOK {
			respData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logging.Print(err)
				return err
			}

			logging.Printf("\n【请求地址】： %s \n【请求参数】： %s \n【请求头】： %s \n【返回】 : %s \n",
				rURL, "上传文件", *header, string(respData))
			return errors.New("上传文件请求成功，上传成功")
		}
		return errors.New("请求失败")
	}
	defer resp.Body.Close()

	return errors.New("请求失败")
}

// 阿里云OSS上传
func AliyunOssUploadFile(endpoint *string, accessKeyId *string, accessKeySecret *string, bucketName *string, fileName string, fileByte *[]byte) string {
	client, err := oss.New(*endpoint, *accessKeyId, *accessKeySecret)
	if err != nil {
		// HandleError(err)
		logging.Printf("创建阿里云上传客户端失败，error：%v", err)
		os.Exit(-1)
	}

	bucket, err := client.Bucket(*bucketName)
	if err != nil {
		// HandleError(err)
		logging.Printf("获取阿里云桶失败，error：%v", err)
		os.Exit(-1)
	}

	err = bucket.PutObject(fileName, bytes.NewReader([]byte(*fileByte)))
	if err != nil {
		// HandleError(err)
		logging.Printf("阿里云上传失败，error：%v", err)
		os.Exit(-1)
	}
	return "https://" + *bucketName + "." + *endpoint + "/" + fileName
}
