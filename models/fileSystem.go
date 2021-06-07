package models

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

//暂时没有用到
type FileSystem struct {
	BucketName   string
	ObjectName   string
	FileLocation string
}

//存储文件信息
type FileInfos struct {
	FileName         string
	FileSize         int64
	FileType         string
	LastModifiedTime time.Time
}

//文件配置
var FileRules map[string]string

//读取文件配置
func InitFileConfig(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	f.Close()
	return config
}
