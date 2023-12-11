package run

import (
	"os"
	"typora_uploader_go/logging"
	"typora_uploader_go/platform"
	"typora_uploader_go/platform/config"
)

func Run(conf *config.PlatformConfig, srcType string, arg *string) string {
	switch srcType {
	case "base64":
		logging.Logger.Printf("暂不支持base64上传")
		os.Exit(-1)
	case "url":
		myhp := platform.NewHttpUploader(*conf)
		myhp.UploadPrepare(arg)
	case "local":
		mylo := platform.NewLocalUploader(*conf)
		mylo.UploadPrepare(arg)
	}

	return ""
}
