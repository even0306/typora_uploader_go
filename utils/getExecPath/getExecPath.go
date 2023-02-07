package getExecPath

import (
	"os"
	"path/filepath"
)

// 获取执行文件当前所在路径
func GetLocalPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}
