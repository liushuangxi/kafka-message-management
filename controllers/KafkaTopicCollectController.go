package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"fmt"

	"github.com/astaxie/beego/orm"

	"github.com/liushuangxi/kafka-message-management/enums"
	"github.com/liushuangxi/kafka-message-management/models"
)

type KafkaTopicCollectController struct {
	BaseController
}

func (c *KafkaTopicCollectController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid", "DataList", "UpdateSeq")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()
}

func (c *KafkaTopicCollectController) Index() {
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl("kafka_topic_collect/index.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "kafka_topic_collect/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "kafka_topic_collect/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("KafkaTopicCollectController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("KafkaTopicCollectController", "Delete")
	c.Data["canAllocate"] = c.checkActionAuthor("KafkaTopicCollectController", "Allocate")
}

func (c *KafkaTopicCollectController) DataGrid() {
	//直接反序化获取json格式的requestbody里的值
	var params models.KafkaTopicCollectQueryParam
	json.Unmarshal(c.Ctx.Input.RequestBody, &params)

	//登陆用户
	params.UserId = int64(c.curUser.Id)

	//获取数据列表和总数
	data, total := models.KafkaTopicCollectPageList(&params)
	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *KafkaTopicCollectController) Delete() {
	strs := c.GetString("ids")
	ids := make([]int, 0, len(strs))
	for _, str := range strings.Split(strs, ",") {
		if id, err := strconv.Atoi(str); err == nil {
			ids = append(ids, id)
		}
	}
	query := orm.NewOrm().QueryTable(models.KafkaTopicCollectTBName())
	if num, err := query.Filter("id__in", ids).Delete(); err == nil {
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("成功删除 %d 项", num), 0)
	} else {
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	}
}

func (c *KafkaTopicCollectController) Collect() {
	userId := c.curUser.Id
	broker := c.GetString("broker")
	topic := c.GetString("topic")

	if c.GetString("collect") == "0" {
		query := orm.NewOrm().QueryTable(models.KafkaTopicCollectTBName())

		query = query.Filter("user_id__iexact", userId)
		query = query.Filter("broker__iexact", broker)
		query = query.Filter("topic__iexact", topic)

		if _, err := query.Delete(); err == nil {
			c.jsonResult(enums.JRCodeSucc, "取消收藏成功", 0)
		} else {
			c.jsonResult(enums.JRCodeFailed, "取消收藏失败", 0)
		}
	} else {
		m := models.KafkaTopicCollect{}

		m.UserId = userId
		m.Broker = broker
		m.Topic = topic

		o := orm.NewOrm()

		if _, err := o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "收藏成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "收藏失败", m.Id)
		}
	}
}