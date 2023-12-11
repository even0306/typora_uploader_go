package run

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"typora_uploader_go/config"
	"typora_uploader_go/upload"
	"typora_uploader_go/utils"
)

type run interface {
	Upload(uploadData *string) *string
}

type MyBase64 struct {
	fmtUrl   string
	upName   string
	filetype string
}

type MyLocal struct {
	fmtUrl   string
	upName   string
	filetype string
}

type MyHttp struct {
	fmtUrl   string
	upName   string
	filetype string
}

func NewBase64Uploader[T *config.Nextcloud | *config.Oss](conf T) *MyBase64 {
	conf.(type)
	return &MyBase64{
		fmtUrl:   ,
		upName:   "",
		filetype: "",
	}
}

func NewLocalUploader[T *config.Nextcloud | *config.Oss](conf T) *MyLocal {
	return &MyLocal{}
}

func NewHttpUploader[T *config.Nextcloud | *config.Oss](conf T) *MyHttp {
	return &MyHttp{}
}

// base64上传
func (b *MyBase64) Upload(args *string) *string {
	user := conf.User
	passwd := conf.Passwd
	b.UploadUrl = conf.Bucket + "/" + user + "/" + conf.Path + "/"
	b.DownloadUrl = conf.Domain + "/" + user + "/" + conf.Path + "/"
	b.Auth = map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+passwd))}
	b.filePath = strings.Split(strings.Split(*args, "base64,")[1], ")")[0]
	file, err := base64.StdEncoding.DecodeString(string(b.filePath))
	if err != nil {
		logging.Printf("解密base64失败，error: %v", err)
	}

	//判断文件格式
	resq.filetype = utils.GetFileExt(&file) //base64如何判断文件类型？？？
	if resq.filetype == "" {
		logging.Printf("文件格式不支持")
		os.Exit(-1)
	}

	b.fileName = utils.CreateUUID() + "." + resq.filetype
	resq.upName = b.UploadUrl + b.fileName
	err = upload.NextcloudUploadFile(resq.upName, &file, &b.Auth)
	if err != nil {
		resq.fmtUrl = b.DownloadUrl + b.fileName + "\n"
	}
	return &resq.fmtUrl
}

// 本地文件上传
func (l *MyLocal) Upload(args *string) *string {
	user := conf.User
	passwd := conf.Passwd
	l.UploadUrl = conf.Bucket + "/" + user + "/" + conf.Path + "/"
	l.DownloadUrl = conf.Domain + "/" + user + "/" + conf.Path + "/"
	l.Auth = map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+passwd))}
	l.filePath = *args
	file, err := utils.ReadFile(&l.filePath)
	if err != nil {
		logging.Printf("读取文件失败，error：%v", err)
	}

	//判断文件格式
	resq.filetype = utils.GetFileExt(file)
	if resq.filetype == "" {
		logging.Printf("文件格式不支持")
		os.Exit(-1)
	}

	l.fileName = utils.CreateUUID() + "." + resq.filetype
	resq.upName = l.UploadUrl + l.fileName
	if conf.PicBed == "nextcloud" {
		err = upload.NextcloudUploadFile(resq.upName, file, &l.Auth)
		if err != nil {
			resq.fmtUrl = l.DownloadUrl + l.fileName + "\n"
		}
	} else if conf.PicBed == "aliyunOss" {
		resq.fmtUrl = upload.AliyunOssUploadFile(&conf.Bucket, &conf.User, &conf.Passwd, &conf.BucketName, l.fileName, file)
	}
	return &resq.fmtUrl
}

// 网络文件上传
func (h *MyHttp) Upload(args *string) *string {
	user := conf.User
	passwd := conf.Passwd
	h.UploadUrl = conf.Bucket + "/" + user + "/" + conf.Path + "/"
	h.DownloadUrl = conf.Domain + "/" + user + "/" + conf.Path + "/"
	h.Auth = map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+passwd))}
	uid := utils.CreateUUID()

	tmp := getexecpath.GetLocalPath() + "/" + uid
	h.filePath = *args
	utils.DownloadFile(h.filePath, tmp)

	file, err := utils.ReadFile(&tmp)
	if err != nil {
		logging.Printf("读取文件失败，error：%v", err)
	}

	//判断文件格式
	resq.filetype = utils.GetFileExt(file)
	if resq.filetype == "" {
		logging.Printf("文件格式不支持")
		os.Exit(-1)
	}
	h.fileName = uid + "." + resq.filetype

	err = os.Remove(tmp)
	if err != nil {
		fmt.Printf("删除缓存图片失败，error：%v", err)
	}
	if conf.PicBed == "nextcloud" {
		err = upload.NextcloudUploadFile(h.UploadUrl+h.fileName, file, &h.Auth)
		if err != nil {
			resq.fmtUrl = h.DownloadUrl + h.fileName + "\n"
		}
	} else if conf.PicBed == "aliyunOss" {
		resq.fmtUrl = upload.AliyunOssUploadFile(&conf.Bucket, &conf.User, &conf.Passwd, &conf.BucketName, h.fileName, file)
	}
	return &resq.fmtUrl
}
