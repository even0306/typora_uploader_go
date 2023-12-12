package config

import (
	"encoding/json"
	"io"
	"os"
	"typora_uploader_go/logging"
)

type PicBed struct {
	Picbed string `json:"picBed"`
}

type Platform struct {
	PicBed          PicBed
	Endpoint        string `json:"endpoint"`
	BucketName      string `json:"bucketName"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	DownloadUrl     string `json:"downloadUrl"`
	UseSSL          bool   `json:"useSSL"`
}

var Config struct {
	ExecPath string
}

func NewReadConfig() *Platform {
	return &Platform{
		PicBed:          PicBed{Picbed: ""},
		Endpoint:        "",
		BucketName:      "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		DownloadUrl:     "",
		UseSSL:          false,
	}
}

func (pf *Platform) ReadConfig(exPath string) {
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

	json.Unmarshal([]byte(byteValue), &pf.PicBed)
	json.Unmarshal([]byte(byteValue), &pf)
}
