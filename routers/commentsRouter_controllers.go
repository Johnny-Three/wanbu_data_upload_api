package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"] = append(beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"] = append(beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"] = append(beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"] = append(beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"] = append(beego.GlobalControllerRouter["wanbu_data_upload_api/controllers:WanbuDataUploadRecordController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

}
