package run

import (
	"typora_uploader_go/config"
	"typora_uploader_go/platform"
)

func Run(conf *config.Platform, srcType string, arg *string) string {
	mylo := platform.NewDataPreparer()
	preConf, fileByte := mylo.UploadPrepare(conf, arg)

	viewURL := platform.UploadSelecter(preConf, fileByte)

	return viewURL
}
