package aliyun_OSS_operation

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

//对象存储客户端，可以使用成员变量client来进行原始操作。
type OSSClient struct {
	client          *oss.Client //操作的客户端
	endPoint        string      //访问域名
	accessKeyID     string      //访问用户的ID
	accessKeySecret string      //访问密钥
}

//定义的OSS客户端在controller层就使用这个东西
var Ossclient OSSClient

//初始化OSS客户端
func init() {
	logs.Info("OSS client initializing......")

	endpoint, _ := web.AppConfig.String("endpoint")
	accesskeyid, _ := web.AppConfig.String("accesskeyid")
	accesskeysecret, _ := web.AppConfig.String("accesskeysecret")
	logs.Info("read endpoint :", endpoint)
	logs.Info("accesskeyid :", accesskeyid)
	logs.Info("accesskeysecret :", accesskeysecret)
	Ossclient = OSSClient{
		//杭州 华东1节点 ID：oss-cn-hangzhou
		endPoint:        endpoint,
		accessKeyID:     accesskeyid,
		accessKeySecret: accesskeysecret,
	}
	//新建客户端 在这里没有填写option 可以查看oss.clientoption中有哪些成员
	var err error
	Ossclient.client, err = oss.New(Ossclient.endPoint, Ossclient.accessKeyID, Ossclient.accessKeySecret)
	//出错就退出了
	if err != nil {
		fmt.Println("Error:", err)
		//os.Exit(-1)
	}
	logs.Info("OSS client OK!")
}

//创建bucket，bucket的命名应该要有规范，需要在前端限制
func (o *OSSClient) CreateBucket(bucketName string) bool {
	isExist, err := o.client.IsBucketExist(bucketName)

	if err != nil {
		handleError(err)
	}

	if isExist {
		logs.Info("bucket aleady exists")
		return false
	}

	err = o.client.CreateBucket(bucketName)
	if err != nil {
		handleError(err)
	}

	return true
}

//上传文件
//params
//bucketName 上传到哪个Bucket
//objectName 上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg
//localFileName 由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
func (o *OSSClient) UploadFile(bucketName string, objectName string, localFileName string) bool {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
		return false
	}

	err = bucket.PutObjectFromFile(objectName, localFileName)

	if err != nil {
		handleError(err)
		return false
	}
	return true
}

//下载文件
//params
//bucketName 下载到哪个Bucket
//objectName 下载文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg
//downloadedFileName 期望下载到本地的路径，由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
func (o *OSSClient) DownloadFile(bucketName string, objectName string, downloadedFileName string) {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}

	err = bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		handleError(err)
	}

}

//列举文件 ls
//params
//bucketName 列举哪个bucket下面的文件
func (o *OSSClient) ListFile(bucketName string) {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 列举文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			handleError(err)
		}
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			fmt.Println("Bucket: ", object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
}

//删除文件 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。
//强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
//params
//bucketName 列举哪个bucket下面的文件
//objectName 删除那个文件？ 表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
func (o *OSSClient) DeleteFile(bucketName string, objectName string) {
	// 获取存储空间。
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 删除文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		handleError(err)
	}
}

//错误处理
func handleError(err error) {
	logs.Info("error encountered：", err.Error())
	os.Exit(-1)
}

//获取存储空间的信息
//https://help.aliyun.com/document_detail/145680.html?spm=a2c4g.11186623.6.1356.604743e1suvwE2
func (o *OSSClient) GetInfo(bucketName string) {
	res, err := o.client.GetBucketInfo(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("BucketInfo.Location: ", res.BucketInfo.Location)
	fmt.Println("BucketInfo.CreationDate: ", res.BucketInfo.CreationDate)
	fmt.Println("BucketInfo.ACL: ", res.BucketInfo.ACL)
	fmt.Println("BucketInfo.Owner: ", res.BucketInfo.Owner)
	fmt.Println("BucketInfo.StorageClass: ", res.BucketInfo.StorageClass)
	fmt.Println("BucketInfo.RedundancyType: ", res.BucketInfo.RedundancyType)
	fmt.Println("BucketInfo.ExtranetEndpoint: ", res.BucketInfo.ExtranetEndpoint)
	fmt.Println("BucketInfo.IntranetEndpoint: ", res.BucketInfo.IntranetEndpoint)
}
