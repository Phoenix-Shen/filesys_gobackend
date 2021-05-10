package controllers

import (
	"FileSys/aliyun_OSS_operation"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

// Operations about FileSystem
type FileSystemController struct {
	beego.Controller
}

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
	print("bucketName是", bucketName, "\n")
	aliyun_OSS_operation.Ossclient.GetInfo(bucketName)
	f.Data["json"] = "succeed"
	f.ServeJSON()
}
