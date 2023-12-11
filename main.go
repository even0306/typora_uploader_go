package main

import (
	"os"
	"path/filepath"
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/platform"
	"typora_uploader_go/utils"
)

func main() {
	//找到执行程序所在位置
	ex, err := os.Executable()
	if err != nil {
		logging.Logger.Panic(err)
	}
	exPath := filepath.Dir(ex)

	PlantformConfig := config.NewReadConfig()
	PlantformConfig.ReadConfig(exPath)

	logging.NewLogger(config.Config.ExecPath + "/server.log")

	for idx, args := range os.Args {
		if idx == 0 {
			continue
		}
		fileType := utils.FileType(&args)

		downloadUrl := platform.PlatformAssert(PlantformConfig, fileType, &args)

		if downloadUrl != "" {
			logging.Logger.Printf("Upload Success:\n")
			logging.Logger.Printf(downloadUrl + "\n")
		}
	}
}
