package aliyun_OSS_operation

import (
	"FileSys/models"
	"errors"
	"fmt"
	"os"
	"path"

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
	//控制台输出密码就是找死
	//logs.Info("read endpoint :", endpoint)
	//logs.Info("accesskeyid :", accesskeyid)
	//logs.Info("accesskeysecret :", accesskeysecret)

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
	//读取文件配置
	models.FileRules = models.InitFileConfig("./conf/rules.conf")
	logs.Info("OSS client OK!")
}

//创建bucket，bucket的命名应该要有规范，需要在前端限制
func (o *OSSClient) CreateBucket(bucketName string) (bool, error) {
	isExist, err := o.client.IsBucketExist(bucketName)

	if err != nil {
		handleError(err)
		return false, err
	}

	if isExist {
		logs.Info("bucket aleady exists")
		return false, errors.New("bucket already exists")
	}

	err = o.client.CreateBucket(bucketName)
	if err != nil {
		handleError(err)
		return false, err
	}

	return true, nil
}

//上传文件
//params
//bucketName 上传到哪个Bucket
//objectName 上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg
//localFileName 由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
func (o *OSSClient) UploadFile(bucketName string, objectName string, localFileName string) (bool, error) {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
		return false, err
	}

	if o.IsExist(bucketName, objectName) {
		handleError(errors.New("duplicate fileName"))
		return false, errors.New("Error:文件名重复")
	}
	forbidOverwirte := oss.ForbidOverWrite(true)
	err = bucket.PutObjectFromFile(objectName, localFileName, forbidOverwirte)

	if err != nil {
		handleError(err)
		return false, err
	}
	return true, nil
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
//返回文件列表
func (o *OSSClient) ListFile(bucketName string) map[string]*models.FileInfos {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 列举文件。
	fileCollection := map[string]*models.FileInfos{}

	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker), oss.MaxKeys(1000))
		if err != nil {
			handleError(err)
			return nil
		}
		var tmp *models.FileInfos
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			fileExt := path.Ext(object.Key)
			if fileExt != ".keep" {
				fileClass, ok := models.FileRules[fileExt]
				if !ok {
					fileClass = fileExt
				}
				tmp = &models.FileInfos{FileName: object.Key, FileSize: object.Size, FileType: fileClass, LastModifiedTime: object.LastModified}
				fileCollection[object.Key] = tmp
			}
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return fileCollection
}

//删除文件 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。
//强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
//params
//bucketName 列举哪个bucket下面的文件
//objectName 删除那个文件？ 表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
func (o *OSSClient) DeleteFile(bucketName string, objectName string) bool {
	// 获取存储空间。
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
		return false
	}
	// 删除文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		handleError(err)
		return false
	}
	return true
}

//错误处理
func handleError(err error) {
	logs.Info("Error encountered：", err.Error())
	//os.Exit(-1)
}

//获取存储空间的信息
//https://help.aliyun.com/document_detail/145680.html?spm=a2c4g.11186623.6.1356.604743e1suvwE2
func (o *OSSClient) GetInfo(bucketName string) map[string]string {
	res, err := o.client.GetBucketInfo(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	/*fmt.Println("BucketInfo.Location: ", res.BucketInfo.Location)
	fmt.Println("BucketInfo.CreationDate: ", res.BucketInfo.CreationDate)
	fmt.Println("BucketInfo.ACL: ", res.BucketInfo.ACL)
	fmt.Println("BucketInfo.Owner: ", res.BucketInfo.Owner)
	fmt.Println("BucketInfo.StorageClass: ", res.BucketInfo.StorageClass)
	fmt.Println("BucketInfo.RedundancyType: ", res.BucketInfo.RedundancyType)
	fmt.Println("BucketInfo.ExtranetEndpoint: ", res.BucketInfo.ExtranetEndpoint)
	fmt.Println("BucketInfo.IntranetEndpoint: ", res.BucketInfo.IntranetEndpoint)*/
	rtval := map[string]string{"BucketInfo.Location": res.BucketInfo.Location,
		"BucketInfo.CreationDate":     res.BucketInfo.CreationDate.Local().String(),
		"BucketInfo.ACL":              res.BucketInfo.ACL,
		"BucketInfo.Owner":            res.BucketInfo.Owner.DisplayName,
		"BucketInfo.StorageClass":     res.BucketInfo.StorageClass,
		"BucketInfo.RedundancyType":   res.BucketInfo.RedundancyType,
		"BucketInfo.ExtranetEndpoint": res.BucketInfo.ExtranetEndpoint,
		"BucketInfo.IntranetEndpoint": res.BucketInfo.IntranetEndpoint}
	return rtval
}

//获取bucket对象
func (o *OSSClient) GetBucket(bucketName string) *oss.Bucket {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	return bucket
}

//查看对象是否存在
func (o *OSSClient) IsExist(bucketName string, fileName string) bool {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	isExist, err := bucket.IsObjectExist(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return isExist
}

//列举指定前缀的存储空间
//以下代码用于列举包含指定前缀（prefix）的存储空间：
func (o *OSSClient) ListPrefix(bucketName string, prefix string) (map[string]*models.FileInfos, error) {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	fileCollection := map[string]*models.FileInfos{}
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker), oss.MaxKeys(1000), oss.Prefix(prefix))
		if err != nil {
			handleError(err)
			return nil, err
		}
		var tmp *models.FileInfos
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			fileExt := path.Ext(object.Key)
			if fileExt != ".keep" {
				fileClass, ok := models.FileRules[fileExt]
				if !ok {
					fileClass = fileExt
				}
				tmp = &models.FileInfos{FileName: object.Key, FileSize: object.Size, FileType: fileClass, LastModifiedTime: object.LastModified}
				fileCollection[object.Key] = tmp
			}

		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return fileCollection, nil
}

//文件重命名
//先拷贝再删除
func (o *OSSClient) RenameFile(bucketName string, oldName string, newName string) error {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	_, err = bucket.CopyObject(oldName, newName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if !o.DeleteFile(bucketName, oldName) {
		return errors.New("delete file failed")
	}
	return nil
}

//列举指定目录下所有子目录的信息
func (o *OSSClient) ListDirInfo(bucketName string, dirName string) ([]string, error) {
	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	dirInfo := []string{}

	marker := oss.Marker("")
	prefix := oss.Prefix(dirName)
	for {
		lor, err := bucket.ListObjects(marker, prefix, oss.Delimiter("/"))
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		/*for _, dirName := range lor.CommonPrefixes {
			dirInfo = append(dirInfo, dirName)
		}*/
		dirInfo = append(dirInfo, lor.CommonPrefixes...)
		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	return dirInfo, nil
}

//创建文件夹
func (o *OSSClient) CreateFolder(bucketName string, dirName string) (bool, error) {
	var err error
	var res bool
	if res, err = o.UploadFile(bucketName, path.Join(dirName, ".keep"), "./.keep"); res {
		return true, nil
	}
	return false, err
}

//删除文件夹
func (o *OSSClient) DeleteFolder(bucketName string, dirName string) (bool, error) {

	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return false, err
	}
	// 列举所有包含指定前缀的文件并删除。
	marker := oss.Marker("")
	prefix := oss.Prefix(dirName)
	count := 0
	for {
		lor, err := bucket.ListObjects(marker, prefix)
		if err != nil {
			fmt.Println("Error:", err)
			return false, err
		}

		objects := []string{}
		for _, object := range lor.Objects {
			objects = append(objects, object.Key)
		}

		delRes, err := bucket.DeleteObjects(objects, oss.DeleteObjectsQuiet(true))
		if err != nil {
			fmt.Println("Error:", err)
			return false, err
		}

		if len(delRes.DeletedObjects) > 0 {
			fmt.Println("these objects deleted failure,", delRes.DeletedObjects)
			return false, errors.New("deleted failure")
		}

		count += len(objects)

		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}
	return true, nil

}
