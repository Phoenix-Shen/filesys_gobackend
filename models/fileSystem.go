package models

import "time"

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
