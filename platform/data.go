package platform

import (
	"os"
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/utils"
)

type dataPrepare interface {
	UploadPrepare(uploadData *string) *string
}

type MyBase64 struct {
	conf        config.Platform
	uploadURL   string
	downloadURL string
	auth        map[string]string
	fmtUrl      string
	upName      string
	filetype    string
}

type MyLocal struct {
	UploadURL       string
	DownloadURL     string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	FileName        string
}

type MyHttp struct {
	conf        config.Platform
	uploadURL   string
	downloadURL string
	auth        map[string]string
	fmtUrl      string
	upName      string
	filetype    string
}

func NewBase64Uploader(conf config.Platform) *MyBase64 {
	return &MyBase64{
		conf:        conf,
		uploadURL:   "",
		downloadURL: "",
		auth:        map[string]string{},
		fmtUrl:      "",
		upName:      "",
		filetype:    "",
	}
}

func NewLocalUploader() *MyLocal {
	return &MyLocal{
		UploadURL:       "",
		DownloadURL:     "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		BucketName:      "",
		FileName:        "",
	}
}

func NewHttpUploader(conf config.Platform) *MyHttp {
	return &MyHttp{
		conf:        conf,
		uploadURL:   "",
		downloadURL: "",
		auth:        map[string]string{},
		fmtUrl:      "",
		upName:      "",
		filetype:    "",
	}
}

// base64上传
// func (b *MyBase64) UploadPrepare(args *string) *string {
// 	accessKeyId := b.conf.AccessKeyId
// 	accessKeySecret := b.conf.AccessKeySecret
// 	b.uploadURL = b.conf.Endpoint + "/" + accessKeyId + "/" + b.conf.BucketName + "/"
// 	b.downloadURL = b.conf.DownloadUrl + "/" + accessKeyId + "/" + b.conf.BucketName + "/"
// 	b.auth = map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(accessKeyId+":"+accessKeySecret))}
// 	b.filePath = strings.Split(strings.Split(*args, "base64,")[1], ")")[0]
// 	file, err := base64.StdEncoding.DecodeString(string(b.filePath))
// 	if err != nil {
// 		logging.Printf("解密base64失败，error: %v", err)
// 	}

// 	//判断文件格式
// 	resq.filetype = utils.GetFileExt(&file) //base64如何判断文件类型？？？
// 	if resq.filetype == "" {
// 		logging.Printf("文件格式不支持")
// 		os.Exit(-1)
// 	}

// 	b.fileName = utils.CreateUUID() + "." + resq.filetype
// 	resq.upName = b.UploadUrl + b.fileName
// 	err = upload.NextcloudUploadFile(resq.upName, &file, &b.Auth)
// 	if err != nil {
// 		resq.fmtUrl = b.DownloadUrl + b.fileName + "\n"
// 	}
// 	return &resq.fmtUrl
// }

// 本地文件上传
func (l *MyLocal) UploadPrepare(conf *config.Platform, args *string) (MyLocal, *[]byte) {
	fileByte, err := utils.ReadFile(args)
	if err != nil {
		logging.Logger.Printf("读取文件失败，error：%v", err)
	}

	//判断文件格式
	filetype := utils.GetFileExt(fileByte)
	if filetype == "" {
		logging.Logger.Printf("文件格式不支持")
		os.Exit(-1)
	}

	fileName := utils.CreateUUID() + "." + filetype
	l.AccessKeyId = conf.AccessKeyId
	l.AccessKeySecret = conf.AccessKeyId

	switch conf.PicBed.Picbed {
	case "nextcloud":
		l.UploadURL = conf.Endpoint + "/" + conf.AccessKeyId + "/" + conf.BucketName + "/" + fileName
		l.DownloadURL = conf.DownloadUrl + "/" + conf.AccessKeyId + "/" + conf.BucketName + "/" + fileName
	case "aliyunOss", "minIO":
		l.UploadURL = conf.Endpoint + "/" + conf.BucketName + "/" + fileName
		l.DownloadURL = conf.BucketName + "." + conf.Endpoint + "/" + fileName
		l.BucketName = conf.BucketName
		l.FileName = fileName
	default:
	}

	return *l, fileByte
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
