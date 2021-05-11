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

var bucketName = "go-api-proj"

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
	aliyun_OSS_operation.Ossclient.GetInfo(bucketName)
	f.Data["json"] = "succeed"
	f.ServeJSON()
}

// @Title UploadFile
// @Description upload a file
// @Param uploadFile fromData multipart.file true "file you want to upload"
// @Success 200 {string} upload succeed
// @Failure 403 file is empty
// @router /uploadfile/ [post]
func (f *FileSystemController) UploadFile() {
	file, h, err := f.GetFile("uploadFile")

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
		f.Data["json"] = ("上传成功！")
	} else {

		f.Data["json"] = ("上传失败！")
	}
	f.ServeJSON()
}
