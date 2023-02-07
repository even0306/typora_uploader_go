package utils

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
)

var fileTypeMap sync.Map

func init() {
	fileTypeMap.Store("ffd8ffe000104a464946", "jpg") //JPEG (jpg)
	fileTypeMap.Store("89504e470d0a1a0a0000", "png") //PNG (png)
	fileTypeMap.Store("47494638396126026f01", "gif") //GIF (gif)
}

func checkType(fileCode *string) (fileType string) {
	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(*fileCode, strings.ToLower(k)) ||
			strings.HasPrefix(k, strings.ToLower(*fileCode)) {
			fileType = v
			return false
		}
		return true
	})
	return
}

// 获取前面结果字节的二进制
func bytesToHexString(src *[]byte) string {
	res := bytes.Buffer{}
	if src == nil || len(*src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range *src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
func GetFileExt(fSrc *[]byte) string {
	fileCode := bytesToHexString(fSrc)
	fileType := checkType(&fileCode)
	return fileType
}
