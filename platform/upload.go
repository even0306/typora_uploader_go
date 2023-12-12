package platform

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
	"typora_uploader_go/logging"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadSelecter(pt MyLocal, b *[]byte) string {
	switch pt.PicBed {
	case "nextcloud":
		viewURL, err := NextcloudUploadFile(pt, b)
		if err != nil {
			logging.Logger.Panicf("未获取到下载地址：%v", err)
		}
		return viewURL
	case "aliyunOss":
		viewURL := AliyunOssUploadFile(pt, b)
		return viewURL
	case "minIO":
		viewURL := MinIOUploadFile(pt, b)
		return viewURL
	default:
		logging.Logger.Println("不支持的平台")
	}
	return ""
}

// 上传接口，传url，文件二进制，参数头
func NextcloudUploadFile(header MyLocal, fileByte *[]byte) (string, error) {
	scheme := "http"
	if header.UseSSL {
		scheme = "https"
	}

	auth := map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(header.AccessKeyId+":"+header.AccessKeySecret))}
	req, err := http.NewRequest("PUT", scheme+"://"+header.UploadURL, bytes.NewBuffer(*fileByte))
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
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, connectTimeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(readWriteTimeout))
				return conn, nil
			},
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
			fmtUrl := scheme + "://" + header.DownloadURL + "\n"
			return fmtUrl, nil
		}
		return "", errors.New("请求失败")
	}
	defer resp.Body.Close()

	return "", errors.New("请求失败")
}

// 阿里云OSS上传
func AliyunOssUploadFile(header MyLocal, fileByte *[]byte) string {
	scheme := "http"
	if header.UseSSL {
		scheme = "https"
	}

	// func AliyunOssUploadFile(endpoint *string, accessKeyId *string, accessKeySecret *string, bucketName *string, fileName string, fileByte *[]byte) string {
	client, err := oss.New(header.UploadURL, header.AccessKeyId, header.AccessKeySecret)
	if err != nil {
		// HandleError(err)
		logging.Logger.Panicf("创建阿里云上传客户端失败，error：%v", err)
	}

	bucket, err := client.Bucket(header.BucketName)
	if err != nil {
		// HandleError(err)
		logging.Logger.Panicf("获取阿里云桶失败，error：%v", err)
	}

	err = bucket.PutObject(header.FileName, bytes.NewReader([]byte(*fileByte)))
	if err != nil {
		// HandleError(err)
		logging.Logger.Panicf("阿里云上传失败，error：%v", err)
	}
	return scheme + "://" + header.DownloadURL
}

// minIO OSS上传
func MinIOUploadFile(header MyLocal, fileByte *[]byte) string {
	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(header.UploadURL, &minio.Options{
		Creds:  credentials.NewStaticV4(header.AccessKeyId, header.AccessKeySecret, ""),
		Secure: header.UseSSL,
	})
	if err != nil {
		logging.Logger.Panicf("初始化minio client失败：%v", err)
	}

	byteReader := bytes.NewReader(*fileByte)

	// Upload the test file with FPutObject
	_, err = minioClient.PutObject(ctx, header.BucketName, header.FileName, byteReader, byteReader.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		logging.Logger.Panicf("发送数据失败：%v", err)
	}

	return header.DownloadURL
}
