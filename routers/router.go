// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"FileSys/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/FileSystem",
			beego.NSInclude(
				&controllers.FileSystemController{},
			),
		),
		beego.NSNamespace("/User",
			beego.NSInclude(
				&controllers.FileSystemController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
