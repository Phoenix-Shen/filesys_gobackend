package controllers

import (
	"FileSys/aliyun_OSS_operation"
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about FileSystem
type FileSystemController struct {
	beego.Controller
}

var defaultbucketName = "go-api-proj"

// @Title CreateBucket
// @Description Create a Bucket
// @Param	body		body 	string	 true		"bucketName for creation"
// @Success 200 create succeed
// @Failure 403 body is empty
// @router / [post]
func (f *FileSystemController) Post() {
	var bucketName string
	json.Unmarshal(f.Ctx.Input.RequestBody, &bucketName)
	result := aliyun_OSS_operation.Ossclient.CreateBucket(bucketName)
	if !result {
		f.Data["json"] = "creation failed"
	} else {
		f.Data["json"] = "creation succeed"
	}
	f.ServeJSON()
}

// @Title GetInfo
// @Description get all information
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 succeed
// @Failure 403 bucket name is empty
// @router /info/:bucketname [get]
func (f *FileSystemController) GetInfo() {
	bucketName := f.GetString(":bucketname")
	//print("bucketName是", bucketName, "\n")
	bucketInfo := aliyun_OSS_operation.Ossclient.GetInfo(bucketName)
	f.Data["json"] = bucketInfo
	f.ServeJSON()
}

// @Title UploadFile
// @Description upload a file
// @Param uploadFile formData multipart.file true "file you want to upload"
// @Param bucketName formData string true "bucketName you want to upload"
// @Success 200 {string} upload succeed
// @Failure 403 file is empty
// @router /uploadfile/ [post]
func (f *FileSystemController) UploadFile() {
	file, h, err := f.GetFile("uploadFile")
	bucketName := f.GetString("bucketName")
	if bucketName == "" {
		logs.Info("检测到您的BucketName为空，自动设置为go-api-proj")
		bucketName = defaultbucketName
	}
	if err != nil {
		logs.Info("文件上传失败：", err.Error())
		f.Data["json"] = ("文件上传失败：" + err.Error())
		return
	}
	FullFileName := path.Base(h.Filename)
	fileExt := path.Ext(h.Filename)
	fileNameWithoutExt := strings.TrimSuffix(FullFileName, fileExt)
	fileName_Time := time.Now().Format("2006-2-15-15-16-10-12")
	defer file.Close()
	f.SaveToFile("uploadFile", "./cache/"+fileNameWithoutExt+"_"+fileName_Time+fileExt)

	if aliyun_OSS_operation.Ossclient.UploadFile(bucketName, FullFileName, "./cache/"+fileNameWithoutExt+"_"+fileName_Time+fileExt) {
		f.Data["json"] = (FullFileName + "上传成功！")
	} else {

		f.Data["json"] = (FullFileName + "上传失败！")
	}
	f.ServeJSON()
}

// @Title GetFilesList
// @Description get files by bucketname
// @Param	bucketName		path 	string	true		"bucketname for "
// @Success 200 {object} list.List
// @Failure 403 :bucketName is empty
// @router /getFileList/:bucketName [get]
func (f *FileSystemController) GetFileList() {
	bucketName := f.GetString(":bucketName")
	//print("bucketName is :", bucketName)
	filemap := aliyun_OSS_operation.Ossclient.ListFile(bucketName)
	f.Data["json"] = filemap
	f.ServeJSON()
}

// @Title DownloadFiles
// @Descreption downloadfiles from speicified bucket
