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

type PicBed struct {
	Picbed string `json:"picBed"`
}

type PlatformConfig struct {
	PicBed          PicBed
	Endpoint        string `json:"bucket"`
	BucketName      string `json:"bucketName"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	DownloadUrl     string `json:"downloadUrl"`
}

var Config struct {
	ExecPath string
}

func NewReadConfig() *PlatformConfig {
	return &PlatformConfig{
		PicBed:          PicBed{},
		Endpoint:        "",
		BucketName:      "",
		AccessKeyId:     "",
		AccessKeySecret: "",
		DownloadUrl:     "",
	}
}

func (pf *PlatformConfig) ReadConfig(exPath string) {
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

	json.Unmarshal([]byte(byteValue), &pf)
}
