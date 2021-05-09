package aliyun_OSS_operation

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type OSSClient struct {
	client          *oss.Client
	endPoint        string
	accessKeyID     string
	accessKeySecret string
}

//定义的OSS客户端在controller层就使用这个东西
var ossclient OSSClient

//初始化OSS客户端
func init() {

}

func (o *OSSClient) CreateBucket(bucketName string) {

}
