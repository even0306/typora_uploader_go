package platform

import (
	"os"
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/run"
)

func PlatformAssert(platformConfig *config.Platform, fileType string, arg *string) string {
	var downloadUrl string

	switch platformConfig.MyPicBed.Picbed {
	case "nextcloud":
		//base64上传存在bug
		if fileType == "base64" {
			// r.url = *run.Run(&bs64, &args)
			logging.Logger.Printf("暂不支持base64上传")
			os.Exit(-1)
		} else if fileType == "url" {
			myHttp := run.NewHttpUploader(platformConfig.MyNextcloud)
			downloadUrl = *myHttp.Upload(arg)
		} else if fileType == "local" {

			downloadUrl = *myLocal.Upload(arg)
		}
	case "aliyunOss", "minIO":
		//base64上传存在bug
		if fileType == "base64" {
			// r.url = *run.Run(&bs64, &args)
			logging.Logger.Printf("暂不支持base64上传")
			os.Exit(-1)
		} else if fileType == "url" {
			downloadUrl = *myhttp.Upload(arg)
		} else if fileType == "local" {
			downloadUrl = *mylocal.Upload(arg)
		}
	default:
		logging.Logger.Panicf("不支持的平台")
	}
	return downloadUrl
}
