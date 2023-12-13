package platform

import (
	"os"
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/utils/assert"

	"github.com/google/uuid"
)

type DataPreparer struct {
	PicBed          string
	UploadURL       string
	DownloadURL     string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	FileName        string
	UseSSL          bool
}

func NewDataPreparer() *DataPreparer {
	return &DataPreparer{
		PicBed:          "",
		UploadURL:       "",
		DownloadURL:     "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "",
		FileName:        "",
		UseSSL:          false,
	}
}

// 准备数据
func (dp *DataPreparer) UploadPrepare(conf *config.Platform, args *string) (DataPreparer, *[]byte) {
	fileByte, err := os.ReadFile(*args)
	if err != nil {
		logging.Logger.Printf("读取文件失败，error：%v", err)
	}

	//判断文件格式
	filetype := assert.GetFileExt(&fileByte)
	if filetype == "" {
		logging.Logger.Panicln("文件格式不支持")
	}

	scheme := "http"
	if conf.UseSSL {
		scheme = "https"
	}

	fileName := uuid.NewString() + "." + filetype
	dp.AccessKeyId = conf.AccessKeyId
	dp.AccessKeySecret = conf.AccessKeySecret
	dp.UseSSL = conf.UseSSL

	switch conf.PicBed.Picbed {
	case "nextcloud":
		dp.PicBed = conf.PicBed.Picbed
		dp.UploadURL = scheme + "://" + conf.Endpoint + "/" + conf.AccessKeyId + "/" + conf.BucketName + "/" + fileName
		dp.DownloadURL = scheme + "://" + conf.DownloadUrl + "/" + conf.AccessKeyId + "/" + conf.BucketName + "/" + fileName
	case "aliyunOss":
		dp.PicBed = conf.PicBed.Picbed
		dp.UploadURL = conf.Endpoint + "/" + conf.BucketName + "/" + fileName
		dp.DownloadURL = conf.BucketName + "." + conf.Endpoint + "/" + fileName
		dp.BucketName = conf.BucketName
		dp.FileName = fileName
	case "minIO":
		dp.PicBed = conf.PicBed.Picbed
		dp.UploadURL = conf.Endpoint
		dp.DownloadURL = scheme + "://" + conf.Endpoint + "/" + conf.BucketName + "/" + fileName
		dp.BucketName = conf.BucketName
		dp.FileName = fileName
	default:
		logging.Logger.Println("不支持的平台")
	}

	return *dp, &fileByte
}

// 网络文件上传
// func (h *MyHttp) UploadPrepare(args *string) *string {
// 	user := conf.User
// 	passwd := conf.Passwd
// 	h.UploadUrl = conf.Bucket + "/" + user + "/" + conf.Path + "/"
// 	h.DownloadUrl = conf.Domain + "/" + user + "/" + conf.Path + "/"
// 	h.Auth = map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+passwd))}
// 	uid := utils.CreateUUID()

// 	tmp := getexecpath.GetLocalPath() + "/" + uid
// 	h.filePath = *args
// 	utils.DownloadFile(h.filePath, tmp)

// 	file, err := utils.ReadFile(&tmp)
// 	if err != nil {
// 		logging.Logger.Printf("读取文件失败，error：%v", err)
// 	}

// 	//判断文件格式
// 	resq.filetype = utils.GetFileExt(file)
// 	if resq.filetype == "" {
// 		logging.Logger.Printf("文件格式不支持")
// 		os.Exit(-1)
// 	}
// 	h.fileName = uid + "." + resq.filetype

// 	err = os.Remove(tmp)
// 	if err != nil {
// 		fmt.Printf("删除缓存图片失败，error：%v", err)
// 	}
// 	if conf.PicBed == "nextcloud" {
// 		err = upload.NextcloudUploadFile(h.UploadUrl+h.fileName, file, &h.Auth)
// 		if err != nil {
// 			resq.fmtUrl = h.DownloadUrl + h.fileName + "\n"
// 		}
// 	} else if conf.PicBed == "aliyunOss" {
// 		resq.fmtUrl = upload.AliyunOssUploadFile(&conf.Bucket, &conf.User, &conf.Passwd, &conf.BucketName, h.fileName, file)
// 	}
// 	return &resq.fmtUrl
// }
