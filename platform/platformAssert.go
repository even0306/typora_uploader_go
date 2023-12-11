package platform

import (
	"os"
	"typora_uploader_go/logging"
	"typora_uploader_go/platform/config"
	"typora_uploader_go/run"
)

func PlatformAssert(platformConfig *config.PlatformConfig, fileType string, arg *string) string {
	var downloadUrl string

	switch fileType {
	case "base64":
		logging.Logger.Printf("暂不支持base64上传")
		os.Exit(-1)
	case "url":
		switch platformConfig.PicBed.Picbed {
		case "nextcloud":

		case "aliyunOss", "minIO":

		default:
			logging.Logger.Panicf("不支持的平台")
		}
	case "local":

	}

	switch platformConfig.MyPicBed.Picbed {
	case "nextcloud":
		//base64上传存在bug
		if fileType == "base64" {
			// r.url = *run.Run(&bs64, &args)
			logging.Logger.Printf("暂不支持base64上传")
			os.Exit(-1)
		} else if fileType == "url" {
			myHttp := run.NewHttpUploader(*platformConfig)
			downloadUrl = *myHttp.Upload(arg)
		} else if fileType == "local" {
			myLocal := run.NewHttpUploader(*platformConfig)
			downloadUrl = *myLocal.Upload(arg)
		}
	case "aliyunOss", "minIO":
		//base64上传存在bug
		if fileType == "base64" {
			// r.url = *run.Run(&bs64, &args)
			logging.Logger.Printf("暂不支持base64上传")
			os.Exit(-1)
		} else if fileType == "url" {
			myHttp := run.NewHttpUploader(*platformConfig)
			downloadUrl = *myHttp.Upload(arg)
		} else if fileType == "local" {
			myLocal := run.NewHttpUploader(*platformConfig)
			downloadUrl = *myLocal.Upload(arg)
		}
	default:
		logging.Logger.Panicf("不支持的平台")
	}
	return downloadUrl
}
