package run

import (
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/platform"
	"typora_uploader_go/upload"
)

func Run(conf *config.Platform, srcType string, arg *string) string {
	mylo := platform.NewLocalUploader()
	preConf, fileByte := mylo.UploadPrepare(conf, arg)
	viewURL, err := upload.NextcloudUploadFile(preConf, fileByte)
	if err != nil {
		logging.Logger.Panicf("未获取到下载地址：%v", err)
	}
	return viewURL
}
