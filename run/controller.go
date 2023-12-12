package run

import (
	"encoding/base64"
	"os"
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/platform"
	"typora_uploader_go/utils"
	"typora_uploader_go/utils/assert"

	"github.com/google/uuid"
)

func Run(conf *config.Platform, arg *string) string {
	if assert.SourceTypeAssert(arg) == "network" {
		uid := uuid.NewString()
		tmp := config.Config.ExecPath + "/." + uid
		utils.CacheFile(*arg, tmp)
		byteTmp, err := os.ReadFile(tmp)
		if err != nil {
			logging.Logger.Panicln(err)
		}
		suf := assert.GetFileExt(&byteTmp)
		tmpFile := tmp + "." + suf
		arg = &tmpFile
		os.Rename(tmp, *arg)
	} else if assert.SourceTypeAssert(arg) == "base" {
		byteTmp, err := base64.StdEncoding.DecodeString(*arg)
		if err != nil {
			logging.Logger.Panicln(err)
		}
		uid := uuid.NewString()
		tmp := config.Config.ExecPath + "/." + uid
		err = os.WriteFile(tmp, byteTmp, 0666)
		if err != nil {
			logging.Logger.Panicln(err)
		}
		suf := assert.GetFileExt(&byteTmp)
		tmpFile := tmp + "." + suf
		arg = &tmpFile
		os.Rename(tmp, *arg)
	}

	datapreparer := platform.NewDataPreparer()
	preConf, fileByte := datapreparer.UploadPrepare(conf, arg)
	viewURL := platform.UploadSelecter(preConf, fileByte)

	err := os.Remove(*arg)
	if err != nil {
		logging.Logger.Println(err)
	}

	return viewURL
}
