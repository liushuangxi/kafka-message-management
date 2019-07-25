package routers

import (
	"github.com/liushuangxi/kafka-message-management/controllers"

	"github.com/astaxie/beego"
)

func init() {

	//KafkaBroker路由
	beego.Router("/kafka_broker/index", &controllers.KafkaBrokerController{}, "*:Index")
	beego.Router("/kafka_broker/datagrid", &controllers.KafkaBrokerController{}, "Get,Post:DataGrid")
	beego.Router("/kafka_broker/edit/?:id", &controllers.KafkaBrokerController{}, "Get,Post:Edit")
	beego.Router("/kafka_broker/delete", &controllers.KafkaBrokerController{}, "Get,Post:Delete")

	//KafkaTopic路由
	beego.Router("/kafka_topic/index", &controllers.KafkaTopicController{}, "*:Index")
	beego.Router("/kafka_topic/datagrid", &controllers.KafkaTopicController{}, "Get,Post:DataGrid")

	//KafkaMessage路由
	beego.Router("/kafka_message/index", &controllers.KafkaMessageController{}, "*:Index")
	beego.Router("/kafka_message/datagrid", &controllers.KafkaMessageController{}, "Get,Post:DataGrid")

	//KafkaTopicCollect路由
	beego.Router("/kafka_topic_collect/index", &controllers.KafkaTopicCollectController{}, "*:Index")
	beego.Router("/kafka_topic_collect/datagrid", &controllers.KafkaTopicCollectController{}, "Get,Post:DataGrid")
	beego.Router("/kafka_topic_collect/delete", &controllers.KafkaTopicCollectController{}, "Post:Delete")
	beego.Router("/kafka_topic_collect/update", &controllers.KafkaTopicCollectController{}, "Post:Update")

	//用户角色路由
	beego.Router("/role/index", &controllers.RoleController{}, "*:Index")
	beego.Router("/role/datagrid", &controllers.RoleController{}, "Get,Post:DataGrid")
	beego.Router("/role/edit/?:id", &controllers.RoleController{}, "Get,Post:Edit")
	beego.Router("/role/delete", &controllers.RoleController{}, "Post:Delete")
	beego.Router("/role/datalist", &controllers.RoleController{}, "Post:DataList")
	beego.Router("/role/allocate", &controllers.RoleController{}, "Post:Allocate")
	beego.Router("/role/updateseq", &controllers.RoleController{}, "Post:UpdateSeq")

	//资源路由
	beego.Router("/resource/index", &controllers.ResourceController{}, "*:Index")
	beego.Router("/resource/treegrid", &controllers.ResourceController{}, "POST:TreeGrid")
	beego.Router("/resource/edit/?:id", &controllers.ResourceController{}, "Get,Post:Edit")
	beego.Router("/resource/parent", &controllers.ResourceController{}, "Post:ParentTreeGrid")
	beego.Router("/resource/delete", &controllers.ResourceController{}, "Post:Delete")

	//快速修改顺序
	beego.Router("/resource/updateseq", &controllers.ResourceController{}, "Post:UpdateSeq")

	//通用选择面板
	beego.Router("/resource/select", &controllers.ResourceController{}, "Get:Select")
	beego.Router("/resource/usermenutree", &controllers.ResourceController{}, "POST:UserMenuTree")
	beego.Router("/resource/checkurlfor", &controllers.ResourceController{}, "POST:CheckUrlFor")

	//后台用户路由
	beego.Router("/backenduser/index", &controllers.BackendUserController{}, "*:Index")
	beego.Router("/backenduser/datagrid", &controllers.BackendUserController{}, "POST:DataGrid")
	beego.Router("/backenduser/edit/?:id", &controllers.BackendUserController{}, "Get,Post:Edit")
	beego.Router("/backenduser/delete", &controllers.BackendUserController{}, "Post:Delete")

	//后台用户中心
	beego.Router("/usercenter/profile", &controllers.UserCenterController{}, "Get:Profile")
	beego.Router("/usercenter/basicinfosave", &controllers.UserCenterController{}, "Post:BasicInfoSave")
	beego.Router("/usercenter/uploadimage", &controllers.UserCenterController{}, "Post:UploadImage")
	beego.Router("/usercenter/passwordsave", &controllers.UserCenterController{}, "Post:PasswordSave")

	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/login", &controllers.HomeController{}, "*:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")
	beego.Router("/home/datareset", &controllers.HomeController{}, "Post:DataReset")

	beego.Router("/home/404", &controllers.HomeController{}, "*:Page404")
	beego.Router("/home/error/?:error", &controllers.HomeController{}, "*:Error")

	beego.Router("/", &controllers.HomeController{}, "*:Index")

}
