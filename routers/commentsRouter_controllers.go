package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"] = append(beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"] = append(beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"],
        beego.ControllerComments{
            Method: "DeleteFile",
            Router: "/delete",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"] = append(beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"],
        beego.ControllerComments{
            Method: "DownloadFiles",
            Router: "/download",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"] = append(beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"],
        beego.ControllerComments{
            Method: "GetFileList",
            Router: "/getFileList/:bucketName",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"] = append(beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"],
        beego.ControllerComments{
            Method: "GetBucketInfo",
            Router: "/info/:bucketname",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"] = append(beego.GlobalControllerRouter["FileSys/controllers:FileSystemController"],
        beego.ControllerComments{
            Method: "UploadFile",
            Router: "/uploadfile/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/:uid",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:uid",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:uid",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: "/login",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["FileSys/controllers:UserController"] = append(beego.GlobalControllerRouter["FileSys/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: "/logout",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
