package models

import "time"

type FileSystem struct {
	BucketName   string
	ObjectName   string
	FileLocation string
}

type FileInfos struct {
	FileName         string
	FileSize         int64
	FileType         string
	LastModifiedTime time.Time
}
