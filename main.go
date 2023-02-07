package main

import (
	"fmt"
	"os"
	"typora_uploader_go/logs"
	"typora_uploader_go/run"
	"typora_uploader_go/utils"
)

func main() {
	var logging = logs.LogFile()
	// var bs64 run.Base64
	var local run.Local
	var http run.Http
	var r struct {
		url string
		req string
	}

	for idx, args := range os.Args {
		if idx == 0 {
			continue
		}
		r.req = utils.FileType(&args)
		//base64上传存在bug
		if r.req == "base64" {
			// r.url = *run.Run(&bs64, &args)
			logging.Printf("暂不支持base64上传")
			os.Exit(-1)
		} else if r.req == "url" {
			r.url = *run.Run(&http, &args)
		} else if r.req == "local" {
			r.url = *run.Run(&local, &args)
		}
		if r.url != "" {
			logging.Printf("Upload Success:\n")
			logging.Printf(r.url + "\n")
			fmt.Printf("Upload Success:\n")
			fmt.Printf(r.url + "\n")
		}
	}
}
