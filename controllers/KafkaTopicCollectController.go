package controllers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

	//Broker List
	brokers := make(map[string]string)

	tmp := make([]*models.KafkaBroker, 0)
	orm.NewOrm().QueryTable(models.KafkaBrokerTBName()).All(&tmp)

	for _, val := range tmp {
		re, _ := regexp.Compile("/$")
		val.Manager = re.ReplaceAllString(val.Manager, "")

		if val.Manager == "" {
			brokers[val.Broker] = ""
		} else {
			brokers[val.Broker] = val.Manager + "/clusters/" + val.Cluster + "/topics/"
		}
	}
	//Broker List

	//Add Kafka Manager
	topics := make([]interface{}, len(data))

	count := 0
	for _, val := range data {
		topic := make(map[string]string)

		topic["Id"] = strconv.Itoa(val.Id)
		topic["Broker"] = val.Broker
		topic["Topic"] = val.Topic
		topic["Alias"] = val.Alias

		_, ok := brokers[val.Broker]

		if ok && brokers[val.Broker] != "" {
			topic["Manager"] = brokers[val.Broker] + val.Topic
		} else {
			topic["Manager"] = ""
		}

		topics[count] = topic
		count++
	}
	//Add Kafka Manager

	//定义返回的数据结构
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = topics
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *KafkaTopicCollectController) Edit() {
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := models.KafkaTopicCollect{Id: Id}
	if Id > 0 {
		o := orm.NewOrm()
		err := o.Read(&m)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	c.Data["m"] = m
	c.setTpl("kafka_topic_collect/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "kafka_topic_collect/edit_footerjs.html"
}

func (c *KafkaTopicCollectController) Save() {
	var err error
	m := models.KafkaTopicCollect{}
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	o := orm.NewOrm()
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}
	}
	c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
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

func (c *KafkaTopicCollectController) Update() {
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
