package config

import (
	"encoding/json"
	"io"
	"os"
	"typora_uploader_go/logging"
)

type config interface {
	ReadConfig(exPath string)
}

type picBed struct {
	Picbed string `json:"picBed"`
}

type Nextcloud struct {
	picBed
	UploadUrl   string `json:"uploadUrl"`
	DownloadUrl string `json:"downloadUrl"`
	Path        string `json:"path"`
	User        string `json:"user"`
	Passwd      string `json:"passwd"`
}

type Oss struct {
	picBed
	Endpoint        string `json:"bucket"`
	BucketName      string `json:"bucketName"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

type PlatformConfig struct {
	MyPicBed    picBed
	MyNextcloud Nextcloud
	MyOss       Oss
}

var Config struct {
	ExecPath string
}

func NewReadConfig() *PlatformConfig {
	return &PlatformConfig{
		MyPicBed:    picBed{},
		MyNextcloud: Nextcloud{},
		MyOss:       Oss{},
	}
}

func (p *PlatformConfig) ReadConfig(exPath string) {
	Config.ExecPath = exPath
	jsonFile, err := os.Open(Config.ExecPath + "/config.json")
	if err != nil {
		logging.Logger.Printf("打开配置文件失败，error：%v", err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		logging.Logger.Printf("读取配置文件失败，error：%v", err)
	}

	json.Unmarshal([]byte(byteValue), &p.MyPicBed)
	json.Unmarshal([]byte(byteValue), &p.MyNextcloud)
	json.Unmarshal([]byte(byteValue), &p.MyOss)
}

// func getConfigValue(byteValue *[]byte, picbed *picBed) {
// 	switch {
// 	case picbed.Picbed == "nextcloud":
// 		var nextcloud Nextcloud
// 		json.Unmarshal([]byte(*byteValue), &nextcloud)
// 		Config.PicBed = nextcloud.Picbed
// 		Config.Bucket = nextcloud.UploadUrl
// 		Config.Domain = nextcloud.DownloadUrl
// 		Config.Path = nextcloud.Path
// 		Config.User = nextcloud.User
// 		Config.Passwd = nextcloud.Passwd
// 	case picbed.Picbed == "aliyunOss":
// 		var aliyunOss aliyunOss
// 		json.Unmarshal([]byte(*byteValue), &aliyunOss)
// 		Config.PicBed = aliyunOss.Picbed
// 		Config.Bucket = aliyunOss.Endpoint
// 		Config.Domain = aliyunOss.Endpoint
// 		Config.BucketName = aliyunOss.BucketName
// 		Config.User = aliyunOss.AccessKeyId
// 		Config.Passwd = aliyunOss.AccessKeySecret
// 	default:
// 		logging.Logger.Print("不支持的图床类型")
// 		os.Exit(-1)
// 	}
// }
