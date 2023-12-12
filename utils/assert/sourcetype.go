package assert

import (
	"strings"
	"typora_uploader_go/logging"
)

// 判断文件来源类型
func SourceTypeAssert(file *string) (tp string) {

	if *file != "" {
		req := strings.Split(*file, ":")[0]
		if req == "data" {
			tp = "base"
		} else if req == "http" || req == "https" {
			tp = "network"
		} else {
			tp = "local"
		}
	} else {
		logging.Logger.Printf("数据不能为空")
	}
	return
}
