package main

import (
	"flag"
	"os"
	"path/filepath"
	"typora_uploader_go/config"
	"typora_uploader_go/logging"
	"typora_uploader_go/run"
	"typora_uploader_go/utils"
)

func main() {
	version := "2.0.0"
	printVersion := flag.Bool("version", false, "[--version]")
	flag.Parse()
	if *printVersion {
		println(version)
		os.Exit(0)
	}

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
		fileType := utils.SrcType(&args)

		downloadUrl := run.Run(PlantformConfig, fileType, &args)

		if downloadUrl != "" {
			logging.Logger.Printf("Upload Success:\n")
			logging.Logger.Printf(downloadUrl + "\n")
		}
	}
}
