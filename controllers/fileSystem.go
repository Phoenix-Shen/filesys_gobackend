package controllers

import (
	"FileSys/aliyun_OSS_operation"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	result, err := aliyun_OSS_operation.Ossclient.CreateBucket(bucketName)
	if !result {
		f.Data["json"] = "Error :" + "creation failed" + err.Error()
	} else {
		f.Data["json"] = "Error :" + "creation succeed" + err.Error()
	}
	f.ServeJSON()
}

// @Title GetBucketInfo
// @Description get all information
// @Param	bucketname		path 	string	true		"the bucket info you want to get"
// @Success 200 succeed
// @Failure 403 bucket name is empty
// @router /info/:bucketname [get]
func (f *FileSystemController) GetBucketInfo() {
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
		logs.Info("Error :"+"文件上传失败：", err.Error())
		f.Data["json"] = ("Error :" + "文件上传失败：" + err.Error())
		f.ServeJSON()
		return
	}
	FullFileName := path.Base(h.Filename)
	fileExt := path.Ext(h.Filename)
	fileNameWithoutExt := strings.TrimSuffix(FullFileName, fileExt)
	fileName_Time := time.Now().Format("2006-2-15-15-16-10-12")
	defer file.Close()
	f.SaveToFile("uploadFile", "./cache/"+fileNameWithoutExt+"_"+fileName_Time+fileExt)

	if succeed, err := aliyun_OSS_operation.Ossclient.UploadFile(bucketName, FullFileName, "./cache/"+fileNameWithoutExt+"_"+fileName_Time+fileExt); succeed {
		f.Data["json"] = (FullFileName + "上传成功！")
	} else {

		f.Data["json"] = ("Error :" + FullFileName + "上传失败！" + err.Error())
	}
	f.ServeJSON()
}

// @Title GetFilesList
// @Description get files by bucketname
// @Param	bucketName		path 	string	true		"bucketname of the filelist"
// @Success 200 list.List
// @Failure 403 :bucketName is empty
// @router /getFileList/:bucketName [get]
func (f *FileSystemController) GetFileList() {
	bucketName := f.GetString(":bucketName")
	//print("bucketName is :", bucketName)
	filemap := aliyun_OSS_operation.Ossclient.ListFile(bucketName)
	//以map形式返回数据
	f.Data["json"] = filemap
	f.ServeJSON()
}

// @Title ListPrefix
// @Description get files by bucketname that maches the prefix
// @Param	bucketName		path 	string	true	"bucketname of the filelist"
// @Param	prefix		    path 	string	true	"prefix"
// @Success 200 list.List
// @Failure 403 :bucketName is empty or prefix is empty
// @router /getFileWithPrefix/:bucketName/:prefix [get]
func (f *FileSystemController) ListPrefix() {
	bucketName := f.GetString(":bucketName")
	prefix := f.GetString(":prefix")
	//print("bucketName is :", bucketName)
	filemap, err := aliyun_OSS_operation.Ossclient.ListPrefix(bucketName, prefix)
	if err != nil {
		f.Data["json"] = err.Error()
		f.ServeJSON()
		return
	}
	//以map形式返回数据
	f.Data["json"] = filemap
	f.ServeJSON()
}

// @Title DownloadFiles
// @Descreption downloadfiles from speicified bucket
// @Param	bucketName		query	string	true		"bucketname for download"
// @Param   fileName        query   string  true        "file Name for download"
// @Success 200 downloadLink
// @Failure 403 :bucketName or filename error
// @router /download [get]
func (f *FileSystemController) DownloadFiles() {
	fileName := f.GetString("fileName")
	bucketName := f.GetString("bucketName")
	//print(fileName, bucketName)
	// 获取存储空间。
	bucket := aliyun_OSS_operation.Ossclient.GetBucket(bucketName)

	// 下载文件到流。
	body, err := bucket.GetObject(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		f.Data["json"] = "Error:" + err.Error()
		f.ServeJSON()
		return
		//os.Exit(-1)
	}

	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Error:", err)
		return
		//os.Exit(-1)
	}
	//下载操作
	//fmt.Println("data:", string(data))
	f.Ctx.ResponseWriter.Header().Add("content-type", "application/octet-stream;charset=utf-8")
	f.Ctx.ResponseWriter.Header().Add("Content-Disposition", "attachement;filename=\""+fileName+"\"")
	f.Ctx.ResponseWriter.Write(data)
}

// @Title DeleteFiles
// @Descreption DeleteFiles from speicified bucket
// @Param	bucketName		query	string	true		"bucketname for delete"
// @Param   fileName        query   string  true        "file Name for delete"
// @Success 200 delete Succeed
// @Failure 403 :bucketName or filename error
// @router /delete [get]
func (f *FileSystemController) DeleteFile() {
	fileName := f.GetString("fileName")
	bucketName := f.GetString("bucketName")
	var result bool
	//判断文件存在否 然后再删除
	if aliyun_OSS_operation.Ossclient.IsExist(bucketName, fileName) {
		result = aliyun_OSS_operation.Ossclient.DeleteFile(bucketName, fileName)
		if result {
			f.Data["json"] = ("succeed :" + fileName + "删除成功！")
		} else {
			f.Data["json"] = ("Error :" + fileName + "删除失败，给老子爬！")
		}
	} else {
		f.Data["json"] = "Error :" + bucketName + "下面的" + fileName + "不存在，给老子爬"
	}

	f.ServeJSON()
}

// @Title RenameFiles
// @Descreption RenameFiles from speicified bucket
// @Param	bucketName	       path	   string  true		   "bucketname for rename"
// @Param   oldfileName        query   string  true        "oldfile Name for rename"
// @Param   newfileName        query   string  true        "newfile Name for rename"
// @Success 200 rename Succeed
// @Failure 403 :bucketName or filenames error
// @router /Rename/:bucketName [get]
func (f *FileSystemController) RenameFiles() {
	oldName := f.GetString("oldfileName")
	newName := f.GetString("newfileName")
	bucketName := f.GetString(":bucketName")
	//print(newName)
	err := aliyun_OSS_operation.Ossclient.RenameFile(bucketName, oldName, newName)
	if err != nil {
		f.Data["json"] = err
	} else {
		f.Data["json"] = "succeed"
	}
	f.ServeJSON()
}
